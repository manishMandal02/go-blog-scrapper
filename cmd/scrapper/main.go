package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/manishmandal02/tech-blog-scrapper/internal/scrapper"
	"github.com/manishmandal02/tech-blog-scrapper/internal/view"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", getHomePageHandler)

	mux.HandleFunc("/scrapper/{blog}", startScrapper)

	// not found handler
	// mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
	// 	http.NotFound(w, r)
	// })

	fmt.Println("🎉 Server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}

func getHomePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		return
	}

	logger.Info("home", r.Method, r.URL.String())

	// articles := scrapper.StartAll()
	path := strings.Split(r.URL.Path, "/")[1]

	view.LayoutComponent(path).Render(r.Context(), w)
}

func startScrapper(w http.ResponseWriter, r *http.Request) {
	// articles := scrapper.StartAll()
	// logger.Info("scrapper", r.Method, r.URL.String())

	blog := r.PathValue("blog")

	fmt.Println("blog:", blog)
	if blog == "" {
		return
	}

	if blog != "all" && blog != "stripe" && blog != "uber" && blog != "netflix" {
		return
	}

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
