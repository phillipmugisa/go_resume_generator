package app

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phillipmugisa/go_resume_generator/data"
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

	contextData := map[string]string{}

	// check is user is logged in

	// if not logged in display landing page
	return RenderHtml(c, w, "landing.html", contextData)
}

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

	// check is user is logged in

	// if not logged in display landing page
	return RenderHtml(c, w, "auth/login.html", contextData)
}

func (a *AppServer) handleSignUp(c context.Context, w http.ResponseWriter, r *http.Request) error {

	contextData := map[string]string{}

	// check is user is logged in redirest

	// if not logged in display landing page

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)

		// make sure password confirmation was successful
		if strings.Compare(r.FormValue("password"), r.FormValue("conform_password")) != 0 {
			contextData["firstname"] = r.FormValue("firstname")
			contextData["lastname"] = r.FormValue("lastname")
			contextData["username"] = r.FormValue("username")
			contextData["email"] = r.FormValue("email")
			contextData["phone"] = r.FormValue("phone")
			contextData["bio_data"] = r.FormValue("bio_data")
			contextData["country"] = r.FormValue("country")
			contextData["start_date"] = r.FormValue("start_date")

			contextData["error_message"] = "Password Mismatch"

			return RenderHtml(c, w, "auth/signup.html", contextData)
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
			return dbWriteErr
		}

		err = a.handleImageUpload(*user, r)
		if err != nil {
			return err
		}

		// send verification link
		http.Redirect(w, r, "/auth/login", http.StatusCreated)
		return nil
	}

	return RenderHtml(c, w, "auth/signup.html", contextData)
}

func (a *AppServer) handleImageUpload(user data.User, r *http.Request) error {
	// read file
	// create destination
	// write file

	// upload images if any
	image, handler, imagesError := r.FormFile("resume_image")
	// check is files were uploaded
	if !errors.Is(imagesError, http.ErrMissingFile) {
		if imagesError != nil {

			// assign default avator
			filedbWriteErr := a.storage.CreateUserimage(user, filepath.Join("media/users/images", "avatar.png"))
			if filedbWriteErr != nil {
				return errors.New("error saving file")
			}
		}
		defer image.Close()

		// save to disk
		destination, fileCreateErr := os.Create(filepath.Join("media/users/images", handler.Filename))
		if fileCreateErr != nil {
			return errors.New("error creating destination file")
		}
		defer destination.Close()

		// Copy the file contents to the destination file
		_, fileWriteErr := io.Copy(destination, image)
		if fileWriteErr != nil {
			return errors.New("error parsing form")
		}

		// save in database
		filedbWriteErr := a.storage.CreateUserimage(user, strings.ReplaceAll(destination.Name(), "\\", "/"))
		if filedbWriteErr != nil {
			os.Remove(destination.Name())
			return errors.New("error saving file")
		}
	}

	return nil
}

func RenderHtml(ctx context.Context, w http.ResponseWriter, template_name string, contextData any) error {

	layout_tmpl := fmt.Sprintf("./Templates/%s", "utils/layout.html")
	template_dir := fmt.Sprintf("./Templates/%s", template_name)

	tmpl, parseError := template.ParseFiles(layout_tmpl, template_dir)

	if parseError != nil {
		return fmt.Errorf("error loading template: %v", parseError)
	}

	// err = tmpl.ExecuteTemplate(w, "layout.html", contextData)
	err := tmpl.Execute(w, contextData)
	if err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}
	return nil
}

func handlerRouterError(c context.Context, e error) {
	fmt.Println(e)
}
