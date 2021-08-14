package service

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type CorsMiddleware struct{
	*mux.Router
	Origin string
}

func (fn CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fn.Origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Page-Size, Page-Index")
	w.Header().Set("Access-Control-Expose-Headers", "Message, Page-Size, Page-Index, Page-Count, Record-Page-Count, Record-Total-Count")

	if strings.ToUpper(r.Method) == "OPTIONS" {
		return
	}
	fn.Router.ServeHTTP(w, r)
}
