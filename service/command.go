package service

import (
	"github.com/nats-io/go-nats"
	"service/service/_base"
	"time"
)

func Command(cmd string, timeout time.Duration, natsUri string, natsOptions []nats.Option, args map[string]interface{}, responseHandler func([]byte)) {
	_base.Command(cmd, timeout, natsUri, natsOptions, args, responseHandler)
}
