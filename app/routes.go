package app

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

type httpHandler func(c context.Context, w http.ResponseWriter, r *http.Request) error

func MakeHTTPHandler(f httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := f(ctx, w, r)
		if err != nil {
			handlerRouterError(ctx, err)
		}
	}
}

func (a *AppServer) handleHomeView(c context.Context, w http.ResponseWriter, r *http.Request) error {
	io.WriteString(w, "home View")
	return nil
}

func handlerRouterError(c context.Context, e error) {
	log.Fatal(e)
}
