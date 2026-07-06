package main

import (
	"fmt"
	"os"

	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/logger"
	"payr/internal/plugins"
	"payr/internal/repository"
	"payr/internal/server"
	"payr/internal/transports"

	"payr/internal/transports/telegram"
)

const (
	DEFAULT_CONFIG_PATH = "./payr.config.json"
)

func main() {
	cmd := cmd.New(cmd.InParams{
		ConfigPath: DEFAULT_CONFIG_PATH,
	})
	params := cmd.Parse()

	switch params.Command {
	case "init":
		cmdInit()
	case "run":
		cmdRun(params.ConfigPath)
	}
}

func cmdInit() {
	config := `{
  "server": {
    "host": "127.0.0.1",
    "port": "8080"
  },

  "plugins": "./plugins",

  "transports": {
    "telegram": {
      "bot_token": "<bot_token>"
    }
  },

  "notifications": {
    "hello": {
      "plugin": "template",
      "settings": {
        "template": "Hello, {{ .Name }}!"
      }
    }
  }
}`

	if err := os.WriteFile(DEFAULT_CONFIG_PATH, []byte(config), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ready! Run: ./payr run")
}

func cmdRun(configPath string) {
	log := logger.New()

	repository := repository.New(repository.Config{
		Path:   configPath,
		Logger: log,
	})
	config := repository.GetAll()

	registry := domain.GetRegistry(config)
	globalSettings := domain.GetGlobalSettings(config)

	pluginsManager := plugins.New(logger.New().WithPackage("plugins"))
	pluginsManager.LoadAll(globalSettings.Plugins)

	transportsManager := transports.New(logger.New().WithPackage("transports"))
	transportsManager.RegisterConstructor("telegram", telegram.New)

	for name, cfg := range registry.Transports {
		constructor := transportsManager.GetConstructor(name)
		transport, err := constructor(log, cfg)
		if err != nil {
			log.Fatal("failed to create transport %s: %v", name, err)
		}

		transportsManager.Register(name, transport)
	}

	srv := server.New(server.Config{
		Logger:            logger.New().WithPackage("server"),
		Host:              globalSettings.Host,
		Port:              globalSettings.Port,
		Registry:          registry,
		PluginsManager:    pluginsManager,
		TransportsManager: transportsManager,
	})

	srv.Start()
}
