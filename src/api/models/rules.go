package models

import "time"

type (
	NotificationRule struct {
		Type              NotificationType
		Name              string
		Description       string
		AmountThreshold   int
		DurationThreshold time.Duration
	}
)
