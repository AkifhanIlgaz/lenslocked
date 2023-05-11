package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filePath string) {
	// Header is just a map.  type Header map[string][]string
	// Set() replaces any existing values for the key, Add() appends to existing values
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// We use filepath.Join() to make our code operating system agnostic
	// While Windows uses "\"(backslach) for path separator other popular operating systems use "/"(forward slash)
	t, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := filepath.Join("templates", "home_go.html")
	executeTemplate(w, templatePath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := filepath.Join("templates", "contact_go.html")
	executeTemplate(w, templatePath)
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
    Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
  </li>
</ul>
`)
}

func main() {
	/*
		"gopls": {
			"ui.SemanticTokens":  true
		}
	*/

	/*
		http.Handler => Interface with the ServeHTTP method
		http.HandlerFunc => A function type that has same arguments as ServeHTTP method. Also implements http.Handler interface
	*/

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "page not found")
		// http.NotFound(w, r)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", r)
}
