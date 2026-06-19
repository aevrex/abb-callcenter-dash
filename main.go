package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
}

func NewApp() *App {
	app := &App{
		router: mux.NewRouter(),
	}

	app.routes()

	return app
}

func (app *App) routes() {
	app.router.HandleFunc("/", handleHome).Methods("GET")
	app.router.HandleFunc("/about", handleAbout).Methods("GET")
}

func (app *App) Run() {
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", app.router))
}

func main() {
	app := NewApp()
	app.Run()
}
