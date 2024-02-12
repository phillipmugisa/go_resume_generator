package app

import (
	"fmt"
	"net/http"
)

type AppServer struct {
	port    int
	storage string
}

func NewAppServer(p int) *AppServer {
	return &AppServer{
		port: p,
	}
}

func (a *AppServer) Run() error {
	sm := http.NewServeMux()
	a.registerRoutes(sm)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: sm,
	}
	return server.ListenAndServe()

}

func (a *AppServer) registerRoutes(sm *http.ServeMux) {
	sm.HandleFunc("/", MakeHTTPHandler(a.handleHomeView))
}
