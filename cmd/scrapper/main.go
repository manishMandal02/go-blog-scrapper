package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/manishmandal02/tech-blog-scrapper/internal/scrapper"
	"github.com/manishmandal02/tech-blog-scrapper/internal/view"
)

type handler struct {
}

type serverHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// loggerMiddleware
	logger.Info("🚀 Server ready")

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// withLogger := func(*http.ServeMux) *http.ServeMux {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		logger.Info(r.Method, r.URL.String())
	// 		h.ServeHTTP(w, r)
	// 	})
	// }

	mux.HandleFunc("/", getHomePageHandler)

	mux.HandleFunc("/all", scrapeAllArticlesHandler)

	fmt.Println("🎉 Server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}

func getHomePageHandler(w http.ResponseWriter, r *http.Request) {

	// articles := scrapper.StartAll()

	view.LayoutComponent().Render(r.Context(), w)
}

func scrapeAllArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articlesJSON, err := os.ReadFile("./internal/scrapper/articles.json")

	if err != nil {
		fmt.Println("error reading articles:", err)
		panic(err)
	}

	articles := []scrapper.Article{}

	json.Unmarshal(articlesJSON, &articles)

	view.ScrapperResult(articles).Render(r.Context(), w)
}

func saveArticlesToFile(articles []scrapper.Article) {
	articlesJSON, err := json.Marshal(articles)

	if err != nil {
		fmt.Println("error marshalling articles:", err)
		panic(err)
	}

	// save articles to file
	os.WriteFile("./internal/scrapper/articles.json", articlesJSON, 0644)
}
