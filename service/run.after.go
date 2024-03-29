package service

import (
	"context"
	"fmt"
	"github.com/go-diary/diary"
	"github.com/gorilla/mux"
	"net/http"
	"service/service/_base"
	"service/service/info"
	"sync"
)

func RunAfter(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	port := fmt.Sprint(info.Args["port"])
	httpsCert := fmt.Sprint(info.Args["httpsCert"])
	httpsKey := fmt.Sprint(info.Args["httpsKey"])
	disableHttps := info.Args["disableHttps"].(bool)
	origin := fmt.Sprint(info.Args["origin"])

	router := mux.NewRouter()

	// serve api html documentation on the root path
	p.Info("http.bind.main", diary.M{
		"path": "/",
	})
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(200)
		_, _ = writer.Write(info.MustAsset("api.html"))
	})

	// serve openapi.json specification file
	p.Info("http.bind.openapi", diary.M{
		"path": "/openapi.json",
	})
	router.HandleFunc("/openapi.json", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(info.MustAsset("openapi.json"))
	})

	// add all requested bindings to web server
	for topic, binding := range _base.Bindings {
		if err := p.Scope("http.bind", func(s diary.IPage) {
			s.Info("data", diary.M{
				"method": binding.Method,
				"path":   binding.Path,
			})
			router.HandleFunc(binding.Path, _base.BindHandler(
				s,
				binding.Timeout,
				_base.TargetLocal(topic),
				binding.Extract,
				binding.ValidateRequest,
				binding.ConvertResponse,
				binding.Permissions...,
			)).Methods(binding.Method)
		}); err != nil {
			panic(err)
		}
	}

	srv := http.Server{
		Addr: ":" + port,
		// the always annoying CORS middleware, for added security of course ;)
		Handler: &_base.CorsMiddleware{Router: router, Origin: origin},
	}
	p.Info("http.server", diary.M{
		"addr": ":" + port,
	})

	// wait for shutdown signal in separate thread
	go func() {
		group.Add(1)
		defer group.Done()

		// closing the shutdown chan will broadcast a close signal
		<-shutdown

		p.Notice("http.server.shutdown", diary.M{
			"addr": ":" + port,
		})

		if err := srv.Shutdown(context.TODO()); err != nil {
			p.Warning("http.server.stop.error", "failed to stop web server", diary.M{
				"addr":     ":" + port,
				"error":    err,
				"errorMsg": err.Error(),
			})
		} else {
			p.Notice("http.server.stop", diary.M{
				"addr": ":" + port,
			})
		}
	}()

	// run web server in separate thread
	go func() {
		group.Add(1)
		defer group.Done()

		p.Notice("http.server.start", diary.M{
			"addr": ":" + port,
		})

		if !disableHttps {
			fmt.Printf("\n\nhttps://127.0.0.1:%s\n\n\n", port)
			if err := srv.ListenAndServeTLS(httpsCert, httpsKey); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		} else {
			fmt.Printf("\n\nhttp://127.0.0.1:%s\n\n\n", port)
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		}
	}()
}
