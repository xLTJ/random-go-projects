package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/waow", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "waow")
	})

	http.ListenAndServe(":8000", r)
}
