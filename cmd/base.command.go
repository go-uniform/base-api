package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"service/service"
)

func Command(name string, handler func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("command:%s", name),
		Short: fmt.Sprintf("Request the running %s to execute the %s command", service.AppName, name),
		Long:  fmt.Sprintf("Request the running %s to execute the %s command", service.AppName, name),
		Run:   handler,
	}
}
