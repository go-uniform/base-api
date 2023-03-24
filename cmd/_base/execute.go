package _base

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"os"
	"service/service/info"
)

var NatsUri string
var NatsCert string
var NatsKey string
var DisableTls bool

var RootCmd = &cobra.Command{
	Use:   info.AppName,
	Short: info.AppDescription,
	Long:  info.AppDescription,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func init() {
	// future: enable support for AWS EventBridge, Azure Event Grid and Google Cloud Pub/Sub support as well to migrate to simplify cloud migration

	// base infrastructure command-line arguments should be included on all cobra commands hence we link it to the RootCmd
	RootCmd.PersistentFlags().StringVarP(&NatsUri, "nats", "n", nats.DefaultURL, "The nats cluster URI")
	RootCmd.PersistentFlags().StringVarP(&NatsCert, "natsCert", "", "/etc/ssl/certs/uniform-nats.crt", "The nats cluster TLS certificate file path")
	RootCmd.PersistentFlags().StringVarP(&NatsKey, "natsKey", "", "/etc/ssl/private/uniform-nats.key", "The nats cluster TLS key file path")
	RootCmd.PersistentFlags().BoolVar(&DisableTls, "disableTls", false, "A flag indicating if service should disable tls encryption")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
