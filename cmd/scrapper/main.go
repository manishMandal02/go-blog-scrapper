package main

import (
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

	mux.HandleFunc("/scrapper/{blog}", handlers.StartScrapper)

	fmt.Println("🎉 Server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}
