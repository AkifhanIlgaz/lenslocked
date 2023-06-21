package main

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/lenslocked/controllers"
	"github.com/AkifhanIlgaz/lenslocked/migrations"
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

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup_go.html", "tailwind_go.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin_go.html", "tailwind_go.html"))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // 32-byte auth key
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", csrfMiddleware(umw.SetUser(r)))
}

/*
	"gopls": {
		"ui.SemanticTokens":  true
	}
*/
