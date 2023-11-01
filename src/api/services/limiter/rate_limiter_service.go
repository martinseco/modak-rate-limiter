package limiter

import (
	"time"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	NotificationRulesDBRepo interface {
		Insert(rule models.NotificationRule) error
		Get(notificationType models.NotificationType) (models.NotificationRule, error)
	}

	UserNotificationsDBRepo interface {
		InsertNotification(userID string, notificationType models.NotificationType, rule models.NotificationRule) error
		GetLastNotifications(userID string, notificationType models.NotificationType) ([]time.Time, error)
	}

	rateLimiterService struct {
		rulesDB NotificationRulesDBRepo
		usersDB UserNotificationsDBRepo
	}
)

func NewRateLimiterService(
	rulesDB NotificationRulesDBRepo,
	usersDB UserNotificationsDBRepo,
) *rateLimiterService {
	return &rateLimiterService{
		rulesDB: rulesDB,
		usersDB: usersDB,
	}
}

func (l *rateLimiterService) Allow(userID string, notificationType models.NotificationType) (bool, error) {
	rule, err := l.rulesDB.Get(notificationType)
	if err != nil {
		return false, err
	}

	lastNotifications, err := l.usersDB.GetLastNotifications(userID, notificationType)
	if err != nil {
		return false, err
	}

	if !isFirstNotification(lastNotifications) && isOverLimit(lastNotifications, rule) {
		return false, nil
	}

	err = l.usersDB.InsertNotification(userID, notificationType, rule)
	if err != nil {
		return false, err
	}

	return true, err
}

func isFirstNotification(lastNotification []time.Time) bool {
	return lastNotification == nil || len(lastNotification) == 0
}

func isOverLimit(lastNotifications []time.Time, rule models.NotificationRule) bool {
	if len(lastNotifications) != 0 && len(lastNotifications) < rule.AmountThreshold {
		return false
	}

	oldestNotification := lastNotifications[0]

	return time.Now().Sub(oldestNotification) < rule.DurationThreshold
}
