package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manishmandal02/tech-blog-scrapper/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handlers.GetHomePageHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		res := make(map[string]string)

		res["status"] = "ok"

		json.NewEncoder(w).Encode(res)

		fmt.Fprintf(w, "{'status':'ok'}")
	})

	mux.HandleFunc("/scrapper/{blog}", handlers.StartScrapper)

	fmt.Println("🎉 Server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}
