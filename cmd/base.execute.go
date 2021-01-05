package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-uniform/base-api/service"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   service.AppName,
	Short: service.AppDescription,
	Long:  service.AppDescription,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
