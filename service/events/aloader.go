package events

import "github.com/go-uniform/uniform"

var MustAsset func(string) []byte
var args uniform.M

func Load(argsMap uniform.M, mustAsset func(string) []byte) {
	// since all events use init method to add themselves this function does nothing
	// calling this function will just the code optimizer from annoyingly removing the "unused" package import
	MustAsset = mustAsset
	args = argsMap
}