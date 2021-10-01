package service

import (
	"service/service/_base"
	"service/service/actions"
	"service/service/commands"
	"service/service/events"
	"service/service/hooks"
)

// load all actions, commands, events and hooks
func init() {
	actions.Load()
	commands.Load()
	events.Load()
	hooks.Load()
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
