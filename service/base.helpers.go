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

// add a topic/handler combination to the actions map
func subscribe(topic string, handler uniform.S) {
	if _, exists := actions[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	actions[topic] = handler
}

// get the index of a string inside of a string array
func indexOf(haystack []string, needle string, caseSensitive bool) int {
	if haystack == nil {
		panic("specify an array to search through")
	}
	if needle == "" {
		panic("specify a string to search for")
	}
	if caseSensitive {
		for i, item := range haystack {
			if item == needle {
				return i
			}
		}
	} else {
		lowNeedle := strings.ToLower(needle)
		for i, item := range haystack {
			if strings.ToLower(item) == lowNeedle {
				return i
			}
		}
	}
	return -1
}

// see if a string array contains a given string
func contains(haystack []string, needle string, caseSensitive bool) bool {
	if haystack == nil {
		panic("specify an array to search through")
	}
	if needle == "" {
		panic("specify a string to search for")
	}
	return indexOf(haystack, needle, caseSensitive) != -1
}

// trim the filterItems from the items array
func filter(items []string, filterItems []string) []string {
	if filterItems == nil || len(filterItems) <= 0 {
		return items
	}
	newItems := make([]string, 0)
	for _, item := range items {
		if contains(filterItems, item, false) {
			continue
		}
		newItems = append(newItems, item)
	}
	return newItems
}
