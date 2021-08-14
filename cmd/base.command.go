package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"service/service"
)

func Command(name string, handler func(cmd *cobra.Command, args []string), description string) *cobra.Command {
	if description == "" {
		description = fmt.Sprintf("Request the running %s to execute the %s command", service.AppName, name)
	}
	return &cobra.Command{
		Use:   fmt.Sprintf("command:%s", name),
		Short: description,
		Long:  description,
		Run:   handler,
	}
}
