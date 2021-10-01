package commands

/* Example
This is just an example command for reference
*/

import (
	"github.com/spf13/cobra"
	"service/cmd/_base"
	"service/service"
	"time"
)

func init() {
	cmd := _base.Command("ping", func(cmd *cobra.Command, args []string) {
		service.Command("ping", time.Second, _base.NatsUri, _base.CompileNatsOptions(), map[string]string{
			// todo: link custom flags to arg values here, example: "custom": custom,
		}, func(data []byte) {
			// todo: handle response data, if any is received
		})
	}, "Ping the currently running service instance")

	// todo: add custom CLI flags here

	_base.RootCmd.AddCommand(cmd)
}