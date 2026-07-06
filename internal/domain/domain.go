package domain

import (
	"encoding/json"

	"payr/internal/repository"
)

type Notification struct {
	Name     string
	Plugin   string
	Settings json.RawMessage
}

type Registry struct {
	Notifications map[string]Notification
	Transports    map[string]json.RawMessage
}

type GlobalSettings struct {
	Host    string
	Port    string
	Plugins string
}

func GetGlobalSettings(registryDTO *repository.Registry) *GlobalSettings {
	return &GlobalSettings{
		Host:    registryDTO.Server.Host,
		Port:    registryDTO.Server.Port,
		Plugins: registryDTO.Plugins,
	}
}

func GetRegistry(registryDTO *repository.Registry) *Registry {
	registry := Registry{
		Notifications: make(map[string]Notification, len(registryDTO.Notifications)),
		Transports:    make(map[string]json.RawMessage, len(registryDTO.Transports)),
	}

	for name, e := range registryDTO.Notifications {
		registry.Notifications[name] = Notification{
			Name:     e.Name,
			Plugin:   e.Plugin,
			Settings: e.Settings,
		}
	}

	for key, t := range registryDTO.Transports {
		registry.Transports[key] = t
	}

	return &registry
}
