package service

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"strings"
)

// use this to separate project specific logic from non-specific logic
func local(topic string) string {
	return fmt.Sprintf("%s.%s", AppProject, strings.TrimPrefix(topic, AppProject + "."))
}

// use this to target a topic for system specific action
func system(topic string) string {
	return local(fmt.Sprintf("action.system.%s", topic))
}

// use this to target a topic for a specific entity item action
func item(entity, action string) string {
	return local(fmt.Sprintf("action.%s.item.%s", entity, action))
}

// use this to target a topic for a specific entity group action
func list(entity, action string) string {
	return local(fmt.Sprintf("action.%s.list.%s", entity, action))
}

// use this to target a topic for a cli level command
func command(topic string) string {
	return local(fmt.Sprintf("command.%s", strings.TrimPrefix(topic, "command.")))
}

// use this to target a service topic for event driven behaviours
func event(service, event string) string {
	return fmt.Sprintf("event.%s.%s", service, event)
}

// use this to target a topic for a function that the service exposes
func routine(function string) string {
	return fmt.Sprintf("%s.%s", AppService, strings.TrimPrefix(function, AppService + "."))
}

func subscribe(topic string, handler uniform.S) {
	if _, exists := actions[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	actions[topic] = handler
}