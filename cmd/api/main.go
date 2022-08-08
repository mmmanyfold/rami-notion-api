package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/mmmanyfold/rami-notion-api/pkg/api"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var port string

	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	notionAPIKey := os.Getenv("NOTION_API_KEY")
	API, err := api.NewAPI(notionAPIKey)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// basic CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// middleware setup
	r.Use(
		corsHandler.Handler,
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,                             // log api request calls
		middleware.StripSlashes,                       // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,                          // recover from panics without crashing server
		middleware.Timeout(3000*time.Millisecond),     // Stop processing after 3 seconds
	)

	// obligatory health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// server static files under public
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	r.Route("/api", func(r chi.Router) {
		r.Route("/sync", func(r chi.Router) {
			r.Get("/db/{name}", API.Sync)
		})
	})

	go scheduleNotionAPISync(port)

	log.Println(fmt.Printf("running on port: %s\n", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		panic(err)
	}
}

func scheduleNotionAPISync(port string) {
	everyHalfHour := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-everyHalfHour.C:
			log.Println("worker::syncing::notion::databases")
			response, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s/api/sync/db/projects", port))
			if err != nil {
				log.Fatalf("failed to retrieve projects from Notion API, err: %v\n", err)
			}
			file, err := os.Create("./public/projects.json")
			if err != nil {
				log.Fatalf("failed to create file for storing projects response from Notion API, err: %v\n", err)
			}

			_, err = io.Copy(file, response.Body)
			if err != nil {
				log.Fatalf("failed to store file for storing projects response from Notion API, err: %v\n", err)
			}
			defer response.Body.Close()
		}
	}
}
