package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lazarosanson/challenge-multithreading-goExpert/internal/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cepHandler := handlers.NewCepHandler()

	r.Get("/ceps/{cep}", cepHandler.GetCep)

	http.ListenAndServe(":8080", r)
}
