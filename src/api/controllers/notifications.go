package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/martinseco/modak-rate-limiter/src/api/errors"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
	"github.com/martinseco/modak-rate-limiter/src/api/utils"
)

type (
	NotificationsService interface {
		SendNotification(request models.NotificationRequest) error
	}

	notificationsController struct {
		notificationsService NotificationsService
	}
)

func NewNotificationsController(s NotificationsService) *notificationsController {
	return &notificationsController{notificationsService: s}
}

func (c *notificationsController) SendNotification(w http.ResponseWriter, r *http.Request) {
	var request models.NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.Handle(w, err)
		return
	}

	err = c.notificationsService.SendNotification(request)
	if err != nil {
		errors.Handle(w, err)
		return
	}

	utils.HandleResponse(w, http.StatusAccepted)
}
