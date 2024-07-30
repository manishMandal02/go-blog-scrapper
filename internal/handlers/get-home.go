package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/manishmandal02/tech-blog-scrapper/internal/view"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func GetHomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	logger.Info("home", r.Method, r.URL.String())

	// articles := scrapper.StartAll()
	path := strings.Split(r.URL.Path, "/")[1]

	view.LayoutComponent(path).Render(r.Context(), w)
}


