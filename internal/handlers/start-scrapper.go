package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/manishmandal02/tech-blog-scrapper/internal/scrapper"
	"github.com/manishmandal02/tech-blog-scrapper/internal/view"
)

func StartScrapper(w http.ResponseWriter, r *http.Request) {
	logger.Info("scrapper", r.Method, r.URL.String())

	blog := r.PathValue("blog")

	fmt.Println("blog:", blog)
	if blog == "" {
		return
	}

	if blog != "all" && blog != "stripe" && blog != "uber" && blog != "netflix" {
		return
	}

	articles := []scrapper.Article{}

	h := r.URL.Query()

	isHeadless := h["headless"][0] == "true"

	fmt.Println("isHeadless:", isHeadless)

	var scrappingErr error

	switch blog {
	case "all":
		articles = append(articles, scrapper.StartAll(isHeadless)...)
	case "stripe":
		stripeArticles, err := scrapper.StripeBlog(-1, isHeadless)
		if err != nil {
			scrappingErr = err
		} else {
			articles = append(articles, stripeArticles...)
		}

	case "uber":
		uberArticles, err := scrapper.UberBlog(-1, isHeadless)
		if err != nil {
			scrappingErr = err
		} else {
			articles = append(articles, uberArticles...)
		}
	case "netflix":
		netflixArticles, err := scrapper.NetflixBlog(-1, isHeadless)
		if err != nil {
			scrappingErr = err
		} else {
			articles = append(articles, netflixArticles...)
		}
	}

	if scrappingErr != nil {
		fmt.Println("Error scrapping blog, error:", scrappingErr)
		http.Error(w, scrappingErr.Error(), http.StatusInternalServerError)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Time.After(articles[j].Time)
	})

	view.ScrapperResult(articles).Render(r.Context(), w)
}
