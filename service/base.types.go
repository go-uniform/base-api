package service

import "encoding/gob"

func init() {
	gob.Register(M{})
}

// A package shorthand for map[string]interface{}
type M map[string]interface{}
