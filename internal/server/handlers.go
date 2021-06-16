package server

import (
	"github.com/golang/glog"
	"net/http"
	"strings"
)

func allowCORS(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}

func openAPI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger/index.html")
}

func registerProxyServer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				allowCORS(w, r)
				return
			}
		}
		if r.Method == "GET" && r.URL.Path == "/api" {
			openAPI(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
