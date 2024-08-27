package handlers

import (
	"net/http"
)

func serveStatic(staticFilePath string) http.Handler {
	fileServer := http.FileServer(http.Dir(staticFilePath))
	return http.StripPrefix("/static/", fileServer)

}
