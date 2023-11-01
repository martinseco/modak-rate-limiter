package models

import "time"

type (
	LastNotifications map[NotificationType][]time.Time
)
