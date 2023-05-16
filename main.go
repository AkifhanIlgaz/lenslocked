package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/AkifhanIlgaz/lenslocked/controllers"
	"github.com/AkifhanIlgaz/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.Parse(filepath.Join("templates", "home_go.html")))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "contact_go.html")))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "faq_go.html")))
	r.Get("/faq", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "about_go.html")))
	r.Get("/about", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "page not found")
		// http.NotFound(w, r)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", r)
}

/*
	"gopls": {
		"ui.SemanticTokens":  true
	}
*/

/*
	http.Handler => Interface with the ServeHTTP method
	http.HandlerFunc => A function type that has same arguments as ServeHTTP method. Also implements http.Handler interface
*/
