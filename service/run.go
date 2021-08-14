package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
	AppService = "service"
	Database = AppProject
)

func Run(p diary.IPage) {
	disableTls := args["disableTls"].(bool)
	port := fmt.Sprint(args["port"])
	tlsCert := fmt.Sprint(args["tlsCert"])
	tlsKey := fmt.Sprint(args["tlsKey"])
	origin := fmt.Sprint(args["origin"])

	// todo: setup logic for common operations like auth
	router := mux.NewRouter()

	// serve api html documentation on the root path
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "text/html")
		writer.WriteHeader(200)
		writer.Write(MustAsset("api.html"))
	})

	// serve openapi.json specification file
	router.HandleFunc("/openapi.json", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(200)
		writer.Write(MustAsset("openapi.json"))
	})

	// serve api javascript client
	router.HandleFunc("/client.js", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(200)
		writer.Write(MustAsset("client.js"))
	})

	// todo: /health endpoint
	// todo: /address/search endpoint
	// todo: /auth/login endpoint
	// todo: /auth/login/code-resend endpoint
	// todo: /auth/login/code-validate endpoint
	// todo: /auth/reset endpoint
	// todo: /auth/reset/{id}/resend endpoint
	// todo: /auth/reset/{id}/validate endpoint
	// todo: /auth/reset/complete endpoint

	// todo: bind all entity endpoints

	if !disableTls {
		if err := http.ListenAndServeTLS(":"+ port, tlsCert, tlsKey, &CorsMiddleware{ router, origin }); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(":"+ port, &CorsMiddleware{ router, origin }); err != nil {
			panic(err)
		}
	}
}