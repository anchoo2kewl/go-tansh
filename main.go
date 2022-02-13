package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
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
	fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:anshuman@biswas.me\">anshuman@biswas.me</a>.")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	_, err := fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
	if err != nil {
		return 
	}
	getRuntime(w, r)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

type Router struct {}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}