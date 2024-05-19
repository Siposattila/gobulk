package bulk

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	g "gorm.io/gorm"
)

func (b *Bulk) unsubscribe(w http.ResponseWriter, r *http.Request) {
	mail := strings.ToLower(r.PathValue("email"))
	if !strings.Contains(mail, "@") {
		http.Error(w, "Invalid parameter!", http.StatusBadRequest)

		return
	}

	console.Normal(mail + " is unsubscribing.")

	var e email.Email
	tx := b.EM.GormORM.First(&e, "email = ? AND status = ?", mail, email.EMAIL_STATUS_ACTIVE)
	if tx.Error != nil && tx.Error == g.ErrRecordNotFound {
		console.Warning("The given email was not found or its not active: " + mail)

		http.Error(w, "Not found the given email "+mail, http.StatusNotFound)

		return
	}

	e.Status = email.EMAIL_STATUS_UNSUBSCRIBED
	b.EM.GormORM.Save(e)

	template, error := template.ParseFiles("unsub.html")
	if error != nil {
		console.Fatal("Failed to load template!")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := template.Execute(w, struct {
		Message     string
		Company     string
		Resubscribe string
	}{
		Message:     "We are very sorry to see you go! 😞",
		Company:     b.Config.CompanyName,
		Resubscribe: b.Config.ResubscribeEndpoint + "/" + e.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Bulk) resubscribe(w http.ResponseWriter, r *http.Request) {
	mail := strings.ToLower(r.PathValue("email"))
	if !strings.Contains(mail, "@") {
		http.Error(w, "Invalid parameter!", http.StatusBadRequest)

		return
	}

	console.Normal(mail + " is resubscribing.")

	var e email.Email
	tx := b.EM.GormORM.First(&e, "email = ? AND status = ?", mail, email.EMAIL_STATUS_UNSUBSCRIBED)
	if tx.Error != nil && tx.Error == g.ErrRecordNotFound {
		console.Warning("The given email was not found or its not unsubscribed: " + mail)

		http.Error(w, "Not found the given email "+mail, http.StatusNotFound)

		return
	}

	e.Status = email.EMAIL_STATUS_ACTIVE
	b.EM.GormORM.Save(e)

	template, error := template.ParseFiles("resub.html")
	if error != nil {
		console.Fatal("Failed to load template!")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := template.Execute(w, struct {
		Message string
		Company string
	}{
		Message: "Welcome back! 🥳",
		Company: b.Config.CompanyName,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Bulk) HttpServer() {
	router := http.NewServeMux()
	router.HandleFunc("GET /unsub/{email}", b.unsubscribe)
	router.HandleFunc("GET /resub/{email}", b.resubscribe)

	server := http.Server{
		Addr:    ":" + b.Config.HttpServerPort,
		Handler: router,
	}

	console.Normal("Http server is listening on port :" + b.Config.HttpServerPort)
	server.ListenAndServe()
}
