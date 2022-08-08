package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:anshuman@biswas.me\">anshuman@biswas.me</a>.")
	if err != nil {
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	_, err := fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
	if err != nil {
		return
	}
	getRuntime(w, r)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
<ul>
  <li>
    <b>Is there a free version?</b>
    Yes! We offer a free trial for 30 days on any paid plans.
  </li>
  <li>
    <b>What are your support hours?</b>
    We have support staff answering emails 24/7, though response
    times may be a bit slower on weekends.
  </li>
  <li>
    <b>How do I contact support?</b>
    Email us - <a href="mailto:anshuman@biswas.me">anshuman@biswas.me</a>
  </li>
</ul>
`)
}

type Router struct{}

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

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	r.Mount("/events", EventsResource{}.Routes())
	r.Mount("/rsvps", RsvpsResource{}.Routes())

	fmt.Println("Starting the server on :3000...")
	err = http.ListenAndServe("0.0.0.0:3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
