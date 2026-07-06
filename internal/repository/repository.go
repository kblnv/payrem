package repository

import (
	"encoding/json"
	"os"

	"payr/internal/logger"
)

type Notification struct {
	Name     string          `json:"name"`
	Plugin   string          `json:"plugin"`
	Settings json.RawMessage `json:"settings"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Registry struct {
	Server        Server                     `json:"server"`
	Plugins       string                     `json:"plugins"`
	Transports    map[string]json.RawMessage `json:"transports"`
	Notifications map[string]Notification    `json:"notifications"`
}

type Config struct {
	Path   string
	Logger *logger.Logger
}

type Repository struct {
	path string
	log  *logger.Logger
}

func New(config Config) *Repository {
	return &Repository{
		path: config.Path,
		log:  config.Logger,
	}
}

func (c *Repository) GetAll() *Registry {
	bytes, err := os.ReadFile(c.path)
	if err != nil {
		c.log.Fatal("failed to read config file: %v", err)
	}

	var registry Registry

	err = json.Unmarshal(bytes, &registry)
	if err != nil {
		c.log.Fatal("failed to parse config: %v", err)
	}

	return &registry
}
