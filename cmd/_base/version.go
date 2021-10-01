package _base

import (
	"fmt"
	"github.com/spf13/cobra"
	service "service/service/_base"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show " + service.AppName + " version information",
	Long:  "Show " + service.AppName + " version information",
	Run:VersionHandler,
}

// we expose the version handler so that each service may modify how their version information is displayed
var VersionHandler = func(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s, build %s\n", service.AppName, service.AppVersion, service.AppCommit)
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
