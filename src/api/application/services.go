package application

import (
	"github.com/martinseco/modak-rate-limiter/src/api/models"
	"time"
)

func buildUsersDB() map[string]models.LastNotifications {
	return make(map[string]models.LastNotifications)
}

func buildRulesDB() map[models.NotificationType]models.NotificationRule {
	rulesDB := make(map[models.NotificationType]models.NotificationRule)

	// populate with sample rules
	rulesDB[models.News] = models.NotificationRule{
		Type:              models.News,
		Name:              "News Rule",
		Description:       "News notifications allow 2 every 30 seconds",
		AmountThreshold:   2,
		DurationThreshold: 30 * time.Second,
	}

	rulesDB[models.Status] = models.NotificationRule{
		Type:              models.Status,
		Name:              "Status Rule",
		Description:       "Status notifications allow 1 every 10 seconds",
		AmountThreshold:   1,
		DurationThreshold: 10 * time.Second,
	}

	rulesDB[models.Marketing] = models.NotificationRule{
		Type:              models.Marketing,
		Name:              "Marketing Rule",
		Description:       "Marketing notifications allow 4 every 60 seconds",
		AmountThreshold:   4,
		DurationThreshold: 60 * time.Second,
	}

	return rulesDB
}
