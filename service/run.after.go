package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/gorilla/mux"
	"net/http"
)

func RunAfter(p diary.IPage) {
	disableTls := args["disableTls"].(bool)
	port := fmt.Sprint(args["port"])
	tlsCert := fmt.Sprint(args["tlsCert"])
	tlsKey := fmt.Sprint(args["tlsKey"])
	origin := fmt.Sprint(args["origin"])

	router := mux.NewRouter()

	// serve api html documentation on the root path
	p.Info("bind.main", diary.M{
		"path": "/",
	})
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(200)
		_, _ = writer.Write(MustAsset("api.html"))
	})

	// serve openapi.json specification file
	p.Info("bind.openapi", diary.M{
		"path": "/openapi.json",
	})
	router.HandleFunc("/openapi.json", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(MustAsset("openapi.json"))
	})

	// serve api javascript client
	p.Info("bind.client", diary.M{
		"path": "/client.js",
	})
	router.HandleFunc("/client.js", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(MustAsset("client.js"))
	})

	for topic, binding := range bindings {
		if err := p.Scope("bind.http", func(s diary.IPage) {
			s.Info("data", diary.M{
				"method": binding.Method,
				"path": binding.Path,
			})
			router.HandleFunc(binding.Path, bindHandler(
				s,
				binding.Timeout,
				topic,
				binding.Extract,
				binding.ValidateRequest,
				binding.ConvertResponse,
				binding.Permissions...
			)).Methods(binding.Method)
		}); err != nil {
			panic(err)
		}
	}

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

	go func() {
		if !disableTls {
			if err := http.ListenAndServeTLS(":"+port, tlsCert, tlsKey, &CorsMiddleware{router, origin}); err != nil {
				panic(err)
			}
		} else {
			if err := http.ListenAndServe(":"+port, &CorsMiddleware{router, origin}); err != nil {
				panic(err)
			}
		}
	}()
}