package service

import (
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"service/service/_base"
)

// wrap args into service layer
var args uniform.M

// wrap base execute to avoid circular reference
func Execute(level string, rate, limit int, test bool, natsUri string, natsOptions []nats.Option, argsMap uniform.M) {
	args = argsMap
	_base.InitializeDiary(test, level, rate)
	_base.Execute(limit, test, natsUri, natsOptions, argsMap, RunBefore, RunAfter)
}