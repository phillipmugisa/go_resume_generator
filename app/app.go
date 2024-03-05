package app

import (
	"fmt"
	"net/http"
	"time"
)

func (a *AppServer) Run() error {
	sm := http.NewServeMux()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", a.port),
		Handler:        sm,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	registerStaticRoutes(sm)
	a.registerRoutes(sm)

	fmt.Printf("Running server on port %s...\n", a.port)
	return server.ListenAndServe()

}

func (a *AppServer) registerRoutes(sm *http.ServeMux) {
	sm.HandleFunc("/", MakeHTTPHandler(a.handleHomeView))
	sm.HandleFunc("/auth/", MakeHTTPHandler(a.handleAuthView))
	sm.HandleFunc("/landing/", MakeHTTPHandler(a.handleLandingView))
}

func registerStaticRoutes(sm *http.ServeMux) {
	staticfileserver := http.FileServer(http.Dir("./static/"))

	sm.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", staticfileserver)
		fs.ServeHTTP(w, r)
	})

	// server media files
	mediafileserver := http.FileServer(http.Dir("./media/"))
	sm.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/media", mediafileserver)
		fs.ServeHTTP(w, r)
	})
}
