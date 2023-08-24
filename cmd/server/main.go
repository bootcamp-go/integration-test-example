package main

import (
	"app/cmd/server/handlers"
	"app/internal/sellers/repository"
	"app/internal/sellers/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// env
	// ...

	// app
	// -> dependencies
	st := storage.NewStorageInMemory(nil)
	rp := repository.NewRepositorySellersDefault(st)
	ct := handlers.NewControllerSellers(rp)

	// -> router
	r := chi.NewRouter()
	r.Get("/sellers/{id}", ct.GetById())
	r.Post("/sellers", ct.Save())

	// -> run
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}