package controllers

import (
	"fmt"
	"net/http"

	"github.com/anchoo2kewl/tansh.us/models"
)

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email            string
		LoggedIn         bool
		IsSignupDisabled bool
	}

	fmt.Println("Here,,,,")
	data.Email = r.FormValue("email")
	data.LoggedIn = false
	data.IsSignupDisabled = false
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Disabled(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email            string
		LoggedIn         bool
		IsSignupDisabled bool
	}
	fmt.Println("Here......,,")
	data.Email = r.FormValue("email")
	data.LoggedIn = false
	data.IsSignupDisabled = true
	u.Templates.New.Execute(w, r, data)
}

type Users struct {
	Templates struct {
		New      Template
		SignIn   Template
		LoggedIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
	RsvpService    *models.RsvpService
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		LoggedIn bool
	}
	data.Email = r.FormValue("email")
	data.LoggedIn = false
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	setCookie(w, CookieUserEmail, data.Email)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Printf("Creating user: %s", email)
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	setCookie(w, CookieUserEmail, email)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	token, err := readCookie(r, CookieSession)
	email, err := readCookie(r, CookieUserEmail)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token, email)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	// fmt.Fprintf(w, "Current user: %s\n", user.Email)

	rsvps, err := u.RsvpService.GetRsvps()

	var data struct {
		Email    string
		LoggedIn bool
		Guests   []models.Rsvps
	}
	data.Email = user.Email
	data.LoggedIn = true
	data.Guests = rsvps.Rsvps

	u.Templates.LoggedIn.Execute(w, r, data)
}

func (u Users) Logout(w http.ResponseWriter, r *http.Request) {

	email, err := readCookie(r, CookieUserEmail)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	u.SessionService.Logout(email)

	deleteCookie(w, CookieSession, "XXXXXX")
	deleteCookie(w, CookieUserEmail, "XXXXXXX")

	http.Redirect(w, r, "/", http.StatusFound)

}
