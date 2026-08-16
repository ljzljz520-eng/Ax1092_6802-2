package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"schoolpaper/internal/workbench"
)

//go:embed web/index.html web/styles.css web/app.js
var webFiles embed.FS

func main() {
	staticFiles, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := workbench.NewFixtureStore()
	handler := workbench.NewHandler(workbench.NewService(store))
	router := handler.Routes(http.FileServer(http.FS(staticFiles)))

	log.Printf("school newspaper workbench listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
