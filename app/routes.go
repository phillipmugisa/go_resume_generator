package app

import (
	"context"
	"net/http"
)

func (a *AppServer) handleLandingView(c context.Context, w http.ResponseWriter, r *http.Request) *HandlerError {

	contextData := map[string]any{}
	// check is user is logged in

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"landing.html"}, contextData)
}

func (a *AppServer) handleHomeView(c context.Context, w http.ResponseWriter, r *http.Request) *HandlerError {

	if r.URL.Path != "/" {
		return &HandlerError{
			code:    http.StatusNotFound,
			message: "address not found",
		}
	}

	// check is user is logged in
	_, err := a.IsAuthenticated(r)
	if err != nil {
		http.Redirect(w, r, "/auth/signin/", http.StatusMovedPermanently)
		return nil
	}

	contextData := map[string]any{}

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"manager/index.html"}, contextData)
}
