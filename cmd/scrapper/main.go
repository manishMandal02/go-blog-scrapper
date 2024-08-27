package main

import (
	"fmt"
	"net/http"

	"github.com/manishmandal02/tech-blog-scrapper/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	routeHandlers := handlers.New("./static")

	mux.Handle("/static/*", routeHandlers.Static)

	mux.HandleFunc("/", routeHandlers.Home)

	mux.HandleFunc("/health", routeHandlers.Health)

	mux.HandleFunc("/scrapper/{blog}", routeHandlers.Scrapper)

	fmt.Println("🎉 Server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("error listening: %v", err)
	}
}
