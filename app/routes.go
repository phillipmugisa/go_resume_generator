package app

import (
	"context"
	"net/http"
)

func (a *AppServer) handleLandingView(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]any{}
	// check is user is logged in

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"landing.html"}, contextData)
}

func (a *AppServer) handleHomeView(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]any{}
	// check is user is logged in
	_, err := a.IsAuthenticated(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
		return nil
	}

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"manager/index.html"}, contextData)
}
