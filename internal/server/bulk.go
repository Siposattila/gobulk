package server

import (
	"net/http"

	"github.com/Siposattila/gobulk/internal/logger"
)

func (s *server) bulk(w http.ResponseWriter, r *http.Request) {
	logger.Debug(r.RemoteAddr)

	//template, error := template.ParseFiles("resub.html")
	//if error != nil {
	//	console.Fatal("Failed to load template!")
	//}

	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//err := template.Execute(w, struct {
	//	Message string
	//	Company string
	//}{
	//	Message: "Welcome back! 🥳",
	//	Company: s.config.GetCompanyName(),
	//})

	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
}
