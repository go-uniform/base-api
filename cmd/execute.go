package cmd

import (
	"service/cmd/_base"

	// load all custom commands
	_ "service/cmd/commands"
)

func Execute() {
	_base.Execute()
}
