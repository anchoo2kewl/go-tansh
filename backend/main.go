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
	"sync"

	"github.com/anchoo2kewl/tansh.us/controllers"
	"github.com/anchoo2kewl/tansh.us/models"
	"github.com/anchoo2kewl/tansh.us/templates"
	"github.com/anchoo2kewl/tansh.us/views"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
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
	r1 := chi.NewRouter()
	r2 := chi.NewRouter()
	r1.Use(middleware.Logger)
	r2.Use(middleware.Logger)

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

	r1.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r2.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r1.Get("/contact", controllers.StaticHandler(
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
		r1.Get("/signup", usersC.Disabled)
	} else {
		fmt.Println("Signups Enabled ...")
		r1.Get("/signup", usersC.New)
		r1.Post("/signup", usersC.Create)
	}

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", "tailwind.gohtml"))

	r1.Get("/signin", usersC.SignIn)
	r1.Post("/signin", usersC.ProcessSignIn)

	usersC.Templates.LoggedIn = views.Must(views.ParseFS(
		templates.FS, "events.gohtml", "tailwind.gohtml"))

	r1.Get("/users/me", usersC.CurrentUser)
	r1.Get("/users/logout", usersC.Logout)

	r2.Mount("/events", EventsResource{}.Routes())
	r2.Mount("/rsvps", RsvpsResource{}.Routes())

	app_port_val, port_present := os.LookupEnv("APP_PORT")
	fmt.Printf("APP_PORT env variable present: %t\n", port_present)

	if !port_present {
		fmt.Println("Using default port of 3000!")
		app_port_val = strconv.Itoa(3000)
	}

	fmt.Println("Starting the server on :{}...", app_port_val)

	internal_port_val, port_present := os.LookupEnv("INTERNAL_PORT")
	fmt.Printf("INTERNAL_PORT env variable present: %t\n", port_present)

	if !port_present {
		fmt.Println("Using default port of 6000!")
		internal_port_val = strconv.Itoa(6000)
	}

	fmt.Println("Starting the internal server on :{}...", internal_port_val)

	isCsrfSecure, err := strconv.ParseBool(os.Getenv("APP_CSRF_SECURE"))

	if isCsrfSecure {
		fmt.Println("CSRF Secure Enabled ...")
	} else {
		fmt.Println("CSRF Secure Disabled ...")
	}

	csrfKey := os.Getenv("APP_CSRF_KEY")
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(isCsrfSecure),
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = http.ListenAndServe("0.0.0.0:"+app_port_val, csrfMw(r1))
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err = http.ListenAndServe("0.0.0.0:"+internal_port_val, r2)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()

}
