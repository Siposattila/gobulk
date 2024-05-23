package server

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/logger"
	"gorm.io/gorm"
)

func (s *server) resubscribe(w http.ResponseWriter, r *http.Request) {
	mail := strings.ToLower(r.PathValue("email"))
	if !email.IsEmail(&mail) {
		http.Error(w, "Invalid parameter!", http.StatusBadRequest)

		return
	}

	logger.LogNormal(r.RemoteAddr + " is resubscribing with " + mail)

	var e email.Email
	tx := s.database.GetEntityManager().GetGormORM().First(&e, "email = ? AND status = ?", mail, interfaces.EMAIL_STATUS_UNSUBSCRIBED)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		logger.LogWarning(mail + " is not found.")
		http.Error(w, "Not found the given email "+mail, http.StatusNotFound)

		return
	}

	e.Status = interfaces.EMAIL_STATUS_ACTIVE
	s.database.GetEntityManager().GetGormORM().Save(e)

	template, error := template.ParseFiles("resub.html")
	if error != nil {
		logger.Fatal("Failed to load template!")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := template.Execute(w, struct {
		Message string
		Company string
	}{
		Message: "Welcome back! 🥳",
		Company: s.config.GetCompanyName(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.LogWarning(mail + " is resubscribed.")
}
