package main

import (
	"backend/internal/repositories"
	"backend/internal/repositories/postgres"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	Domain       string
	DSN          string
	DB           repositories.Repository
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	TMDBAPIKey   string
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies timezone=UTC connect_timeout=5", "Postgress connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "bigandrandomsecret", "JWT secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "JWT issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "JWT audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "Cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "Application domain")
	flag.StringVar(&app.TMDBAPIKey, "tmdb-api-key", "9a27a3af220c836a5b5323277c9e53e6", "API key for The Movies DB")
	flag.Parse()

	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal(err)
	}

	app.DB = &postgres.PostgresRepository{DB: conn}

	defer app.DB.GetConnection().Close()

	app.auth = Auth{
		Secret:        app.JWTSecret,
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting application on port ", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())

	if err != nil {
		log.Fatal(err)
	}
}
