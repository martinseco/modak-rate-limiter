package models

type NotificationType string

const (
	News      NotificationType = "news"
	Status    NotificationType = "status"
	Marketing NotificationType = "marketing"
)

type (
	NotificationRequest struct {
		UserID  string           `json:"user_id,omitempty"`
		Type    NotificationType `json:"type,omitempty"`
		Message string           `json:"message,omitempty"`
	}
)
