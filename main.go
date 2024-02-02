package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
)

type application struct {
	session  *sessions.Session
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		errorLog.Fatal(err)
	}

	secret := os.Getenv("SECRET_KEY")

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		session:  session,
	}

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler: app.routes(),
		Addr:    port,
	}

	app.infoLog.Printf("Server starting on port %v", port)
	err := srv.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
