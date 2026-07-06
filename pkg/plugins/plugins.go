package plugins

import "encoding/json"

type Context struct {
	NotificationMeta json.RawMessage
}

type Plugin interface {
	Execute(context *Context) (string, error)
}

type Constructor func(rawConfig json.RawMessage) (Plugin, error)
