package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/anchoo2kewl/tansh.us/controllers"
	"github.com/anchoo2kewl/tansh.us/models"
	"github.com/anchoo2kewl/tansh.us/templates"
	"github.com/anchoo2kewl/tansh.us/views"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func getRuntime(w http.ResponseWriter, r *http.Request) {
	myOS, myArch := runtime.GOOS, runtime.GOARCH
	inWSL := "outside"
	cmd := exec.Command("uname", "-a")
	if output, err := cmd.Output(); err == nil {
		if strings.Contains(strings.ToLower(string(output)), "microsoft") {
			inWSL = "inside"
		}
	}
	_, _ = fmt.Fprintf(w, "Hello, %s!\n", r.UserAgent())
	_, _ = fmt.Fprintf(w, "I'm running on %s/%s.\n", myOS, myArch)
	_, _ = fmt.Fprintf(w, "I'm running %s of WSL.\n", inWSL)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	dbUser, dbPassword, dbName, dbHost :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST")

	database, err := Initialize(dbUser, dbPassword, dbName, dbHost)

	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer database.Conn.Close()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	userService := models.UserService{
		DB: DB,
	}

	sessionService := models.SessionService{
		DB: DB,
	}

	rsvpService := models.RsvpService{
		DB: DB,
	}

	// Setup our controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
		RsvpService:    &rsvpService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", "tailwind.gohtml"))

	isSignupDisabled, err := strconv.ParseBool(os.Getenv("APP_DISABLE_SIGNUP"))

	if isSignupDisabled {
		fmt.Println("Signups Disabled ...")
		r.Get("/signup", usersC.Disabled)
	} else {
		fmt.Println("Signups Enabled ...")
		r.Get("/signup", usersC.New)
		r.Post("/signup", usersC.Create)
	}

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", "tailwind.gohtml"))

	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	usersC.Templates.LoggedIn = views.Must(views.ParseFS(
		templates.FS, "events.gohtml", "tailwind.gohtml"))

	r.Get("/users/me", usersC.CurrentUser)
	r.Get("/users/logout", usersC.Logout)

	r.Mount("/events", EventsResource{}.Routes())
	r.Mount("/rsvps", RsvpsResource{}.Routes())

	port_val, port_present := os.LookupEnv("APP_PORT")
	fmt.Printf("APP_PORT env variable present: %t\n", port_present)

	if !port_present {
		fmt.Println("Using default port of 3000!")
		port_val = strconv.Itoa(3000)
	}

	fmt.Println("Starting the server on :{}...", port_val)

	// csrfKey := os.Getenv("APP_CSRF_KEY")
	// csrfMw := csrf.Protect(
	// 	[]byte(csrfKey),
	// 	// TODO: Fix this before deploying
	// 	csrf.Secure(false),
	// )
	err = http.ListenAndServe("0.0.0.0:"+port_val, r)
	if err != nil {
		log.Fatal(err)
	}
}
