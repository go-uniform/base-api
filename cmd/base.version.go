package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"service/service"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show " + service.AppName + " version information",
	Long:  "Show " + service.AppName + " version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s, build %s\n", service.AppName, service.AppVersion, service.AppCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
