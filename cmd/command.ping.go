package cmd

import (
	"github.com/spf13/cobra"
	"service/service"
)

func init() {
	rootCmd.AddCommand(Command("ping", func(cmd *cobra.Command, args []string) {
		service.InitializeDiary(test, level, rate)
		service.Command("ping", natsUri, compileNatsOptions(), map[string]string{})
	}))
}