package service

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"strings"
)

func local(topic string) string {
	return fmt.Sprintf(AppProject + ".%s", strings.TrimPrefix(topic, AppProject + "."))
}

func command(topic string) string {
	return fmt.Sprintf("command.%s", strings.TrimPrefix(topic, "command."))
}

func subscribe(topic string, handler uniform.S) {
	if _, exists := actions[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	actions[topic] = handler
}