package handlers

import "net/http"

type httpHandler func(w http.ResponseWriter, r *http.Request)

type Handlers struct {
	Static   http.Handler
	Home     httpHandler
	Scrapper httpHandler
	Health   httpHandler
}

func New(staticFilePath string) *Handlers {
	return &Handlers{
		Static:   serveStatic(staticFilePath),
		Home:     getIndexPage,
		Scrapper: startScrapper,
		Health:   health,
	}
}


