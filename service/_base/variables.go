package _base

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

// A private package level variable that contains the service's run-time arguments
var args uniform.M

// A private package level variable that contains the running diary instance
var d diary.IDiary

// A private package level variable that contains the running diary instance's sample trace rate value
var traceRate int

// A private package level variable that contains the running uniform event bus connection instance
var c uniform.IConn

// A private package level variable that indicates if service is in test mode
var testMode = false

// A private package level variable that contains all subscribing actions
var actions = make(map[string]uniform.S)

// A private package level variable that contains all topic subscriptions
var subscriptions = make(map[string]uniform.ISubscription)