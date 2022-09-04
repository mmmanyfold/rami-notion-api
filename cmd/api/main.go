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
	port := getPort()
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

	go scheduleNotionAPISync()
	if os.Getenv("MODE") == "production" {
		go initialNotionDBSync()
	}

	log.Println(fmt.Printf("running on port: %s\n", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		panic(err)
	}
}

func scheduleNotionAPISync() {
	everyHalfHour := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-everyHalfHour.C:
			log.Println("syncing notion databases")

			if err := fetchDBAndPersist("projects"); err != nil {
				log.Println(err)
			}
			time.Sleep(30 * time.Second)
			if err := fetchDBAndPersist("cv-additional"); err != nil {
				log.Println(err)
			}
			time.Sleep(30 * time.Second)
			if err := fetchDBAndPersist("cv-exhibitions-and-screenings"); err != nil {
				log.Println(err)
			}
			time.Sleep(30 * time.Second)
			if err := fetchDBAndPersist("info"); err != nil {
				log.Println(err)
			}
		}
	}
}

func initialNotionDBSync() {
	log.Println("initial databases sync")

	if err := fetchDBAndPersist("projects"); err != nil {
		log.Println(err)
	}
	time.Sleep(30 * time.Second)
	if err := fetchDBAndPersist("cv-additional"); err != nil {
		log.Println(err)
	}
	time.Sleep(30 * time.Second)
	if err := fetchDBAndPersist("cv-exhibitions-and-screenings"); err != nil {
		log.Println(err)
	}
	time.Sleep(30 * time.Second)
	if err := fetchDBAndPersist("info"); err != nil {
		log.Println(err)
	}
}

func fetchDBAndPersist(dbname string) error {
	port := getPort()
	response, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s/api/sync/db/%s", port, dbname))
	if err != nil {
		return fmt.Errorf("error gettingfailed to retrieve projects from Notion API, err: %v\n", err)
	}
	file, err := os.Create(fmt.Sprintf("./public/%s.json", dbname))
	if err != nil {
		return fmt.Errorf("failed to create file for storing %s response from Notion API, err: %v\n", err, dbname)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to store file for storing %s response from Notion API, err: %v\n", err, dbname)
	}

	defer response.Body.Close()
	return nil
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}
