package service

import (
	"service/service/_base"
	"service/service/actions"
	"service/service/commands"
	"service/service/events"
	"service/service/hooks"
)

var MustAsset = _base.MustAsset

// load all actions, commands, events and hooks
func init() {
	actions.Load(args, _base.MustAsset)
	commands.Load(args, _base.MustAsset)
	events.Load(args, _base.MustAsset)
	hooks.Load(args, _base.MustAsset)
	_base.AppClient = AppClient
	_base.AppProject = AppProject
	_base.AppService = AppService
	_base.Database = Database
}

// wrap base const to avoid circular reference
const (
	AppName        = _base.AppName
	AppDescription = _base.AppDescription
	AppVersion     = _base.AppVersion
	AppCommit      = _base.AppCommit
)

// wrap base type to avoid circular reference
type M _base.M
