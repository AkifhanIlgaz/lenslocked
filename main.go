package main

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/lenslocked/controllers"
	"github.com/AkifhanIlgaz/lenslocked/models"
	"github.com/AkifhanIlgaz/lenslocked/templates"
	"github.com/AkifhanIlgaz/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home_go.html", "tailwind_go.html"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact_go.html", "tailwind_go.html"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq_go.html", "tailwind_go.html"))
	r.Get("/faq", controllers.FAQ(tpl))

	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup_go.html", "tailwind_go.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin_go.html", "tailwind_go.html"))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // 32-byte auth key
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", csrfMiddleware(r))
}

// func ExerciseMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ip := r.RemoteAddr
// 		path := r.URL.Path
// 		start := time.Now()
// 		h.ServeHTTP(w, r)

// 		log.Printf("%v made a request to %v. %v", ip, path, time.Since(start))
// 	})
// }

/*
	"gopls": {
		"ui.SemanticTokens":  true
	}
*/

/*
	http.Handler => Interface with the ServeHTTP method
	http.HandlerFunc => A function type that has same arguments as ServeHTTP method. Also implements http.Handler interface
*/
