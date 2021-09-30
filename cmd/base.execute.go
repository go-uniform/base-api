package cmd

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"os"
	"service/service"
)

var natsUri string
var natsCert string
var natsKey string
var level string
var rate int
var limit int
var test bool
var disableTls bool

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

func init() {
	rootCmd.PersistentFlags().StringVarP(&natsUri, "nats", "n", nats.DefaultURL, "The nats cluster URI")
	rootCmd.PersistentFlags().StringVarP(&natsCert, "nats-cert", "", "/etc/ssl/certs/ssl-bundle.crt", "The nats cluster TLS certificate file path")
	rootCmd.PersistentFlags().StringVarP(&natsKey, "nats-key", "", "/etc/ssl/private/ssl.key", "The nats cluster TLS key file path")
	rootCmd.PersistentFlags().BoolVar(&disableTls, "disable-tls", false, "A flag indicating if service should disable tls encryption")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
