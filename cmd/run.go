package cmd

import (
	"github.com/spf13/cobra"
	"service/service"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run " + service.AppName + " service",
	Long:  "Run " + service.AppName + " service",
	Run: func(cmd *cobra.Command, args []string) {
		service.InitializeDiary(test, level, rate)
		service.Execute(limit, test, natsUri, compileNatsOptions(), service.M{
			// todo: add custom args here
		})
	},
}

func init() {
	runCmd.Flags().StringVarP(&level, "lvl", "l", "trace", "The logging level ['trace', 'debug', 'info', 'notice', 'warning', 'error', 'fatal'] that service is running in")
	runCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "The sample rate of the trace logs used for performance auditing [set to -1 to log every trace]")
	runCmd.Flags().IntVarP(&limit, "limit", "x", 1000, "The messages per second that each topic worker will be limited to [set to 0 or less for maximum throughput]")
	runCmd.Flags().BoolVar(&test, "test", false, "A flag indicating if service should enter into test mode")
	runCmd.Flags().BoolVar(&disableTls, "disable-tls", false, "A flag indicating if service should disable tls encryption")

	rootCmd.AddCommand(runCmd)
}
