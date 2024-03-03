package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/phillipmugisa/go_resume_generator/data"
)

func (a *AppServer) handleAuthView(c context.Context, w http.ResponseWriter, r *http.Request) error {

	url_parts := strings.Split(r.URL.String(), "auth")
	switch url_parts[1] {
	case "/login/":
		return a.handleLogin(c, w, r)
	case "/signup/":
		return a.handleSignUp(c, w, r)
	default:
		return errors.New("page not found")
	}
}

func (a *AppServer) handleLogin(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]string{}

	if r.Method == http.MethodPost {
		// login user in
		r.ParseForm()

		contextData["username"] = r.FormValue("username")

		username := r.FormValue("username")
		password := r.FormValue("password")

		err := a.Login(username, password, w)
		if err != nil {
			contextData["error_message"] = "Invalid username/password."
			return a.RenderHtml(c, w, r, []string{"auth/login.html"}, contextData)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return nil
	}

	// check is user is logged in
	_, err := a.IsAuthenticated(r)
	if err == nil {
		// user is already logged in
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return nil
	}

	// if not logged in display landing page
	return a.RenderHtml(c, w, r, []string{"auth/login.html"}, contextData)
}

func (a *AppServer) handleSignUp(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]string{}

	// check is user is logged in redirest

	// if not logged in display landing page

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)

		contextData["firstname"] = r.FormValue("firstname")
		contextData["lastname"] = r.FormValue("lastname")
		contextData["username"] = r.FormValue("username")
		contextData["email"] = r.FormValue("email")
		contextData["phone"] = r.FormValue("phone")
		contextData["bio_data"] = r.FormValue("bio_data")
		contextData["country"] = r.FormValue("country")
		contextData["start_date"] = r.FormValue("start_date")

		// make sure password confirmation was successful
		if strings.Compare(r.FormValue("password"), r.FormValue("conform_password")) != 0 {

			contextData["error_message"] = "Password Mismatch"

			return a.RenderHtml(c, w, r, []string{"auth/signup.html"}, contextData)
		}

		user, err := data.NewUser(
			r.FormValue("firstname"),
			r.FormValue("lastname"),
			r.FormValue("username"),
			r.FormValue("email"),
			r.FormValue("password"),
			r.FormValue("phone"),
			r.FormValue("bio_data"),
			r.FormValue("country"),
			r.FormValue("start_date"),
		)
		if err != nil {
			return err
		}

		// save user data to database
		dbWriteErr := a.storage.CreateUser(*user)
		if dbWriteErr != nil {
			contextData["error_message"] = "Username/Email Not available."
			return a.RenderHtml(c, w, r, []string{"auth/signup.html"}, contextData)
		}

		err = a.HandleImageUpload(*user, r)
		if err != nil {
			return err
		}

		// send verification link

		contextData["success_message"] = "Activation link sent to your email"
		return a.RenderHtml(c, w, r, []string{"auth/login.html"}, contextData)
	}

	return a.RenderHtml(c, w, r, []string{"auth/signup.html"}, contextData)
}
