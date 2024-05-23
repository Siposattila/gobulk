package server

import (
	"net/http"
	"text/template"

	"github.com/Siposattila/gobulk/internal/logger"
)

func (s *server) bulk(w http.ResponseWriter, r *http.Request) {
	logger.LogNormal(r.RemoteAddr + " is accessing the bulk page!")

	template, error := template.ParseFiles("bulk.html")
	if error != nil {
		logger.Fatal("Failed to load template!")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := template.Execute(w, struct{}{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
