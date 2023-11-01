package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BuildRouter(app *application) http.Handler {
	fmt.Println("Building router...")
	mux := chi.NewRouter()

	mux.Get("/ping", pong)

	// Notifications
	mux.Post("/notifications", app.NotificationsController.SendNotification)

	return mux
}

func pong(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "pong...")
}
