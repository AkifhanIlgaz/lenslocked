package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/AkifhanIlgaz/lenslocked/controllers"
	"github.com/AkifhanIlgaz/lenslocked/migrations"
	"github.com/AkifhanIlgaz/lenslocked/models"
	"github.com/AkifhanIlgaz/lenslocked/templates"
	"github.com/AkifhanIlgaz/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DBNAME"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil

}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(cfg)
	if err != nil {
		panic(err)
	}
	// Setup the database

}

func run(cfg config) error {
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}

	// Setup the services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	passwordResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup middlewares
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		EmailService:         emailService,
		PasswordResetService: passwordResetService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup_go.html", "tailwind_go.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin_go.html", "tailwind_go.html"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-password_go.html", "tailwind_go.html"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email_go.html", "tailwind_go.html",
	))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-password_go.html", "tailwind_go.html",
	))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries-new_go.html", "tailwind_go.html",
	))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries-edit_go.html", "tailwind_go.html",
	))
	galleriesC.Templates.Index = views.Must(views.ParseFS(
		templates.FS,
		"galleries-index_go.html", "tailwind_go.html",
	))
	galleriesC.Templates.Show = views.Must(views.ParseFS(
		templates.FS,
		"galleries-show_go.html", "tailwind_go.html",
	))

	// Setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "home_go.html", "tailwind_go.html"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact_go.html", "tailwind_go.html"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq_go.html", "tailwind_go.html"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Get("/users/me", usersC.CurrentUser)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
			r.Post("/{id}/images", galleriesC.UploadImage)
		})
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
