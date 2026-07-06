package server

import (
	"encoding/json"
	"net/http"

	"payr/internal/domain"
	"payr/internal/logger"
	"payr/internal/plugins"
	"payr/internal/transports"

	api "payr/pkg/plugins"
)

type Server struct {
	server            *http.Server
	registry          *domain.Registry
	pluginsManager    *plugins.Manager
	transportsManager *transports.Transports
	log               *logger.Logger
}

type Config struct {
	Logger            *logger.Logger
	Host              string
	Port              string
	Registry          *domain.Registry
	PluginsManager    *plugins.Manager
	TransportsManager *transports.Transports
}

type RequestTransport struct {
	Transport string `json:"name"`
	To        string `json:"to"`
}

type NotificationTriggerRequestBody struct {
	Notification string             `json:"notification"`
	Meta         json.RawMessage    `json:"meta"`
	Transports   []RequestTransport `json:"transports"`
}

func New(config Config) *Server {
	mux := http.NewServeMux()

	server := Server{
		server: &http.Server{
			Addr:    config.Host + ":" + config.Port,
			Handler: mux,
		},
		registry:          config.Registry,
		pluginsManager:    config.PluginsManager,
		transportsManager: config.TransportsManager,
		log:               config.Logger,
	}

	mux.HandleFunc("/notification", server.handleNotificationTrigger)

	return &server
}

func (s *Server) Start() {
	s.log.Info("server is listening on %v", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		s.log.Fatal("server error: %v", err)
	}
}

func (s *Server) handleNotificationTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var payload NotificationTriggerRequestBody
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.log.Warn("failed to parse request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	s.log.Info("handling notification: %v", payload.Notification)

	notification, ok := s.registry.Notifications[payload.Notification]
	if !ok {
		s.log.Warn("notification not found: %v", payload.Notification)
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}

	constructor := s.pluginsManager.GetConstructor(notification.Plugin)
	if constructor == nil {
		s.log.Error("plugin constructor not found: %v", notification.Plugin)
		http.Error(w, "plugin not found", http.StatusInternalServerError)
		return
	}

	plugin, err := constructor(notification.Settings)
	if err != nil {
		s.log.Error("failed to create plugin %s: %v", notification.Plugin, err)
		http.Error(w, "failed to create plugin", http.StatusInternalServerError)
		return
	}

	ctx := api.Context{
		NotificationMeta: payload.Meta,
	}

	result, err := plugin.Execute(&ctx)
	if err != nil {
		s.log.Error("plugin execution failed: %v", err)
		http.Error(w, "plugin execution failed", http.StatusInternalServerError)
		return
	}

	if result == "" {
		s.log.Debug("plugin result: <empty message>")
		w.Write([]byte("ok"))
		return
	}

	s.log.Debug("plugin result: %v", result)

	if len(payload.Transports) == 0 {
		s.log.Warn("no transports specified in request")
		http.Error(w, "no transports specified in request", http.StatusInternalServerError)
		return
	}

	for _, trigger := range payload.Transports {
		transport := s.transportsManager.Get(trigger.Transport)
		if transport == nil {
			s.log.Error("transport not found: %v", trigger.Transport)
			http.Error(w, "transport not found", http.StatusInternalServerError)
			return
		}

		err := transport.Send(result, trigger.To)
		if err != nil {
			s.log.Error("transport send failed: %v", err)
			http.Error(w, "transport send failed", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("ok"))
}
