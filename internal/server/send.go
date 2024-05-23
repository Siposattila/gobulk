package server

import (
	"net/http"

	"github.com/Siposattila/gobulk/internal/logger"
)

func (s *server) send(w http.ResponseWriter, r *http.Request) {
	logger.LogNormal(r.RemoteAddr + " is sending bulk emails!")

	r.ParseForm()
	subject := r.FormValue("subject")
	greeting := r.FormValue("greeting")
	message := r.FormValue("message")
	farewell := r.FormValue("farewell")
	shouldContinue := r.FormValue("shouldContinue") == "on"

	go s.app.GetBulk().Start(subject, greeting, message, farewell, shouldContinue)

	http.Redirect(w, r, "/bulk", http.StatusMovedPermanently)
}
