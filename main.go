package main

import (
	"log"
	"net/http"
	"shorters/controller"
	"shorters/repository"
	"shorters/service"
	"shorters/service/jwt"
	"shorters/service/mail"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	linkController := controller.NewLinkController(service.NewLinkService(repository.NewSQLLinkRepository()))
	authController := controller.NewAuthController(jwt.NewJWTService(), mail.NewMailService())

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(authController.AuthenticateMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "view/index.html")
	})
	router.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "view/signin.html")
	})
	router.Get("/verify", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "view/verify.html")
	})
	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})

	router.Get("/{key}", linkController.Redirect)

	router.Post("/shorten", linkController.Shorten)
	router.Post("/custom-shorten", linkController.CustomShorten)
	router.Post("/signin", authController.SignIn)
	router.Post("/verify", authController.Verify)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
