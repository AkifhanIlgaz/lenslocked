package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/AkifhanIlgaz/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filePath string) {

	t, err := views.Parse(filePath)

	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
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
	templatePath := filepath.Join("templates", "faq_go.html")
	executeTemplate(w, templatePath)
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
