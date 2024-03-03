package app

import (
	"context"
	"net/http"
)

func (a *AppServer) handleHomeView(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]any{}

	// list last 10 users
	users, err := a.storage.GetUsers(map[string]string{})
	if err != nil {
		return err
	}

	contextData["users"] = users

	// check is user is logged in

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"landing.html", "./partials/_userlist.html"}, contextData)
}
