package cmd

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"go-uniform/base-api/service"
)

var natsUri string
var level string
var rate int
var test bool

var database string
var auth string
var port int

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run " + service.AppName + " service",
	Long:  "Run " + service.AppName + " service",
	Run: func(cmd *cobra.Command, args []string) {
		handler := diary.HumanReadableHandler
		if test {
			// test mode is used for creating an automated testing environment
			fmt.Println("entering test mode")
			handler = nil
		}
		service.Execute(test, natsUri, database, auth, level, rate, port, handler, service.M{})
	},
}

func init() {
	runCmd.Flags().StringVarP(&natsUri, "nats", "n", nats.DefaultURL, "The nats cluster URI")
	runCmd.Flags().StringVarP(&level, "lvl", "l", "trace", "The logging level that service is running in")
	runCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "The sample rate of the trace logs used for performance auditing [set to -1 to log every trace]")
	runCmd.Flags().BoolVar(&test, "test", false, "A flag indicating if service should enter into test mode")

	runCmd.Flags().StringVarP(&database, "database", "d", "", "The database environment to use")
	runCmd.Flags().StringVarP(&auth, "auth", "a", "<<database>>", "The auth environment to use")
	runCmd.Flags().IntVarP(&rate, "port", "p", 8000, "The port to serve the API on")

	if err := runCmd.MarkFlagRequired("database"); err != nil {
		panic(err)
	}
	if err := runCmd.MarkFlagRequired("auth"); err != nil {
		panic(err)
	}
	if err := runCmd.MarkFlagRequired("port"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(runCmd)
}
