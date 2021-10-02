package service

import (
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"service/service/_base"
	"service/service/info"

	// load all actions, commands, events, hooks, entities and integrations
	_ "service/service/actions"
	_ "service/service/commands"
	_ "service/service/events"
	_ "service/service/hooks"
	_ "service/service/entities"
	_ "service/service/integrations"

	// load bindings (base-api only)
	_ "service/service/binds"
)

// wrap base execute to avoid circular reference
func Execute(level string, rate, limit int, test bool, natsUri string, natsOptions []nats.Option, argsMap uniform.M) {
	info.Args = argsMap
	if info.Args == nil {
		info.Args = uniform.M{}
	}
	_base.InitializeDiary(test, level, rate)
	_base.Execute(limit, test, natsUri, natsOptions, RunBefore, RunAfter)
}