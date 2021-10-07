package _base

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"service/service/info"
	"strings"
)

// use this to separate project specific logic from non-specific logic
var TargetLocal = func(topic string) string {
	if topic == "" {
		return info.AppProject
	}
	return fmt.Sprintf("%s.%s", info.AppProject, strings.TrimPrefix(topic, info.AppProject+ "."))
}

// use this to target a topic for system specific action
var TargetSystem = func(topic string) string {
	if topic == "" {
		return "system"
	}
	return fmt.Sprintf("system.%s", strings.TrimPrefix(topic, "system."))
}

// use this to target a topic for a specific entity item action
var TargetItem = func(entity, action string) string {
	if entity == "" {
		return TargetLocal("item")
	}
	if action == "" {
		return TargetLocal(fmt.Sprintf("item.%s", entity))
	}
	return TargetLocal(fmt.Sprintf("item.%s.%s", entity, action))
}

// use this to target a topic for a specific entity group action
var TargetList = func(entity, action string) string {
	if entity == "" {
		return TargetLocal("list")
	}
	if action == "" {
		return TargetLocal(fmt.Sprintf("list.%s", entity))
	}
	return TargetLocal(fmt.Sprintf("list.%s.%s", entity, action))
}

// use this to target a topic for a cli level command
var TargetCommand = func(topic string) string {
	return TargetLocal(fmt.Sprintf("command.%s", strings.TrimPrefix(topic, "command.")))
}

// use this to target a service topic for event driven behaviours
var TargetEvent = func(service, event string) string {
	if service == "" {
		return "event"
	}
	if event == "" {
		return fmt.Sprintf("event.%s", service)
	}
	return fmt.Sprintf("event.%s.%s", service, event)
}

// use this to target a topic for a function that the service exposes
var TargetAction = func(service, action string) string {
	if service == "" {
		return "action"
	}
	if action == "" {
		return fmt.Sprintf("action.%s", service)
	}
	return fmt.Sprintf("action.%s.%s", service, action)
}

// add a topic/handler combination to the actions map
var Subscribe = func(topic string, handler uniform.S) {
	if _, exists := actions[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	actions[topic] = handler
}
