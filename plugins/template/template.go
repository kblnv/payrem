package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"payr/pkg/plugins"
)

type Config struct {
	Template string `json:"template"`
}

type TemplatePlugin struct {
	template *template.Template
}

func New(rawConfig json.RawMessage) (plugins.Plugin, error) {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	tmpl, err := template.New("message").Parse(config.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &TemplatePlugin{
		template: tmpl,
	}, nil
}

func (t *TemplatePlugin) Execute(context *plugins.Context) (string, error) {
	var data map[string]any

	if len(context.NotificationMeta) > 0 {
		if err := json.Unmarshal(context.NotificationMeta, &data); err != nil {
			return "", err
		}
	} else {
		data = map[string]any{}
	}

	var buf bytes.Buffer

	if err := t.template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
