package app

import (
	"fmt"
	"net/http"
)

func (a *AppServer) Run() error {
	sm := http.NewServeMux()

	registerStaticRoutes(sm)
	a.registerRoutes(sm)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.port),
		Handler: sm,
	}

	fmt.Printf("Running server on port %s...\n", a.port)
	return server.ListenAndServe()

}

func (a *AppServer) registerRoutes(sm *http.ServeMux) {
	sm.HandleFunc("/", MakeHTTPHandler(a.handleHomeView))
	sm.HandleFunc("/auth/", MakeHTTPHandler(a.handleAuthView))
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
