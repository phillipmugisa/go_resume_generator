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
	"golang.org/x/crypto/bcrypt"
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

func (a *AppServer) HandleImageUpload(user data.User, r *http.Request) error {
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

func (a *AppServer) RenderHtml(ctx context.Context, w http.ResponseWriter, r *http.Request, templates []string, contextData any) error {
	// pass user data to template by default if user is authenticated

	var template_dirs []string

	layout_tmpl := fmt.Sprintf("./Templates/%s", "utils/layout.html")

	template_dirs = append(template_dirs, layout_tmpl)
	for _, t := range templates {
		template_dir := fmt.Sprintf("./Templates/%s", t)
		template_dirs = append(template_dirs, template_dir)
	}

	tmpl, parseError := template.ParseFiles(template_dirs...)
	if parseError != nil {
		return fmt.Errorf("error loading template: %v", parseError)
	}

	// tmpl = tmpl.Funcs(template.FuncMap{
	// 	"getLoggedInUser": func() data.User {
	// 		user, auth_err := a.IsAuthenticated(r)
	// 		if auth_err == nil {
	// 			return *user
	// 		}
	// 		return data.User{}
	// 	},
	// })

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

func checkSessionKey(key string) error {
	// TODO
	// depending on the set session restrictions
	if key == "" {
		return errors.New("invalid session Key")
	}

	return nil
}

func (a AppServer) IsAuthenticated(r *http.Request) (*data.User, error) {
	// get sessionid cookie
	cookie, err := r.Cookie("session_key")
	if err != nil {
		return nil, err
	}

	// asssure key value is present and not edited
	if err := checkSessionKey(cookie.Value); err != nil {
		return nil, err
	}

	// session_key was found, check database for same key
	session, err := a.storage.GetSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	// is the session expired
	if session.Expired {
		err := a.storage.DeleteSession(*session)
		if err != nil {
			return nil, err
		}

		// user should log in again
		return nil, errors.New("expired session")
	}

	if session.Expires_on.Equal(time.Now()) || time.Now().After(session.Expires_on) {
		err := a.storage.CancelSession(*session)
		if err != nil {
			return nil, err
		}

		// user should log in again
		return nil, errors.New("expired session")
	}

	// get the user
	users, err := a.storage.GetUsers(map[string]string{"id": session.User.Id})
	if err != nil {
		return nil, errors.New("error getting user")
	}

	return users[0], nil
}

func (a AppServer) Login(username, password string, w http.ResponseWriter) error {
	// get user if the same username
	users, err := a.storage.GetUsers(map[string]string{"username": fmt.Sprint(username)})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("invalid username/password")
	}
	// compare password
	pwd_err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
	if pwd_err != nil {
		return err
	}

	// valid credentials provided, log in the user
	// create new session
	session, err := users[0].NewSession()
	if err != nil {
		return err
	}

	ss_creation_err := a.storage.CreateSession(*session)
	if ss_creation_err != nil {
		return ss_creation_err
	}

	cookie := http.Cookie{
		Name:     "session_key",
		Value:    session.Key,
		Path:     "/",
		MaxAge:   int(data.Session_duration * 3600),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	return nil
}

func (a AppServer) Logout(w http.ResponseWriter, r *http.Request) error {
	// check is user is logged in
	_, err := a.IsAuthenticated(r)
	if err != nil {
		return err
	}

	// get sessionid cookie
	cookie, err := r.Cookie("session_key")
	if err != nil {
		return err
	}

	// session_key was found, check database for same key
	session, err := a.storage.GetSession(cookie.Value)
	if err != nil {
		return err
	}

	ss_cancel_err := a.storage.CancelSession(*session)
	if ss_cancel_err != nil {
		return ss_cancel_err
	}

	// remove session key from header
	new_cookie := http.Cookie{
		Name:    "session_key",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &new_cookie)

	return nil
}
