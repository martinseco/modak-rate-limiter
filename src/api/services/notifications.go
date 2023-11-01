package services

import (
	"fmt"

	"github.com/martinseco/modak-rate-limiter/src/api/errors"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	NotificationsRateLimiter interface {
		Allow(userID string, notificationType models.NotificationType) (bool, error)
	}

	NotificationsSender interface {
		Send(request models.NotificationRequest) error
	}

	notificationsService struct {
		rateLimiter NotificationsRateLimiter
		sender      NotificationsSender
	}
)

func NewNotificationsService(
	rl NotificationsRateLimiter,
	s NotificationsSender,
) *notificationsService {
	return &notificationsService{
		rateLimiter: rl,
		sender:      s,
	}
}

// SendNotification sends notification to the user. If notification type reaches limiting rules, returns error.
func (s *notificationsService) SendNotification(request models.NotificationRequest) error {
	//check rate limit
	allowed, err := s.rateLimiter.Allow(request.UserID, request.Type)
	if err != nil {
		return err
	}

	if !allowed {
		return errors.ForbiddenError(
			fmt.Sprintf("%s's notification type for %s user reached it's limit!", request.Type, request.UserID),
		)
	}

	//send notification
	err = s.sender.Send(request)
	if err != nil {
		return err
	}

	return nil
}
