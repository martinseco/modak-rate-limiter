package application

import (
	"fmt"
	"net/http"

	"github.com/martinseco/modak-rate-limiter/src/api/controllers"
	"github.com/martinseco/modak-rate-limiter/src/api/repository"
	"github.com/martinseco/modak-rate-limiter/src/api/services"
	"github.com/martinseco/modak-rate-limiter/src/api/services/limiter"
	"github.com/martinseco/modak-rate-limiter/src/api/services/sender"
)

type (
	NotificationsController interface {
		SendNotification(w http.ResponseWriter, r *http.Request)
	}

	application struct {
		NotificationsController NotificationsController
	}
)

func BuildApplication() *application {
	fmt.Println("Building Application...")

	rulesDBRepo := repository.NewRulesDBRepo(buildRulesDB())

	usersDBRepo := repository.NewUsersDBRepo(buildUsersDB())

	rateLimiterService := limiter.NewRateLimiterService(rulesDBRepo, usersDBRepo)

	senderService := sender.NewSenderService()

	notificationsService := services.NewNotificationsService(rateLimiterService, senderService)

	notificationsController := controllers.NewNotificationsController(notificationsService)

	return &application{
		NotificationsController: notificationsController,
	}
}
