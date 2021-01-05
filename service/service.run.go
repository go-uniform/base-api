package service

import (
	"github.com/go-diary/diary"
	"strings"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
)

func Run(p diary.IPage) {
	// todo: set your config values based on the given environment
	switch strings.ToLower(env) {
	default:
		panic("unknown environment given")
	case "production":
	case "prod":
		BaseApiUrl = "https://api.uniform.tech"
		break
	case "demo":
		BaseApiUrl = "https://demo-api.uniform.tech"
		break
	case "qa":
	case "staging":
		BaseApiUrl = "https://qa-api.uniform.tech"
		break
	case "development":
	case "dev":
		BaseApiUrl = "https://dev-api.uniform.tech"
		break
	case "localhost":
	case "local":
		BaseApiUrl = "http://localhost:9000"
		break
	}

	// todo: run required database migrations and/or index creation routines
}