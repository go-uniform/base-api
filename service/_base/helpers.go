package _base

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"strings"
)

// use this to separate project specific logic from non-specific logic
var TargetLocal = func(topic string) string {
	return fmt.Sprintf("%s.%s", AppProject, strings.TrimPrefix(topic, AppProject+ "."))
}

// use this to target a topic for system specific action
var TargetSystem = func(topic string) string {
	return TargetLocal(fmt.Sprintf("action.system.%s", topic))
}

// use this to target a topic for a specific entity item action
var TargetItem = func(entity, action string) string {
	return TargetLocal(fmt.Sprintf("action.%s.item.%s", entity, action))
}

// use this to target a topic for a specific entity group action
var TargetList = func(entity, action string) string {
	return TargetLocal(fmt.Sprintf("action.%s.list.%s", entity, action))
}

// use this to target a topic for a cli level command
var TargetCommand = func(topic string) string {
	return TargetLocal(fmt.Sprintf("command.%s", strings.TrimPrefix(topic, "command.")))
}

// use this to target a service topic for event driven behaviours
var TargetEvent = func(service, event string) string {
	return fmt.Sprintf("event.%s.%s", service, event)
}

// use this to target a topic for a function that the service exposes
var TargetRoutine = func(function string) string {
	return fmt.Sprintf("%s.%s", AppService, strings.TrimPrefix(function, AppService+ "."))
}

// add a topic/handler combination to the actions map
var Subscribe = func(topic string, handler uniform.S) {
	if _, exists := actions[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	actions[topic] = handler
}
