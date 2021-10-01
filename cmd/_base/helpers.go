package _base

import (
	"crypto/tls"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	service "service/service/_base"
)

func CompileNatsOptions() []nats.Option {
	var natsOptions = make([]nats.Option, 0)
	if !DisableTls {
		cert, err := tls.LoadX509KeyPair(NatsCert, NatsKey)
		if err != nil {
			panic(err)
		}
		config := &tls.Config{
			// since NATS backbone should always be on a private line with self-signed certs, we just skip host verification
			InsecureSkipVerify: true,

			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		natsOptions = append(natsOptions, nats.Secure(config))
	}
	return natsOptions
}

func Command(name string, handler func(cmd *cobra.Command, args []string), description string) *cobra.Command {
	if description == "" {
		description = fmt.Sprintf("Request the running %s to execute the %s command", service.AppName, name)
	}
	return &cobra.Command{
		Use:   fmt.Sprintf("command:%s", name),
		Short: description,
		Long:  description,
		Run:   handler,
	}
}
