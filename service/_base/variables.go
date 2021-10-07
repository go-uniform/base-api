package _base

import (
	"github.com/go-uniform/uniform"
)

// A private package level variable that contains all subscribing actions
var actions = make(map[string]uniform.S)

// A private package level variable that contains all topic subscriptions
var subscriptions = make(map[string]uniform.ISubscription)