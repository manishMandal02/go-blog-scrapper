package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	res := make(map[string]string)

	res["status"] = "ok"

	json.NewEncoder(w).Encode(res)

	fmt.Fprintf(w, "{'status':'ok'}")
}
