package repository

import (
	"sync"
	"time"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	usersDBRepo struct {
		db   map[string]models.LastNotifications
		lock sync.Mutex
	}
)

func NewUsersDBRepo(db map[string]models.LastNotifications) *usersDBRepo {
	return &usersDBRepo{
		db:   db,
		lock: sync.Mutex{},
	}
}

func (r *usersDBRepo) InsertNotification(
	userID string,
	notificationType models.NotificationType,
	rule models.NotificationRule,
) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	notifications, found := r.db[userID]
	if found {
		lastNotifications := notifications[notificationType]
		if notificationsLessThanThreshold(lastNotifications, rule.AmountThreshold) {
			//remove oldest
			lastNotifications = lastNotifications[1:]
		}

		//append newest
		lastNotifications = append(lastNotifications, time.Now())
		notifications[notificationType] = lastNotifications
	} else {
		notifications = map[models.NotificationType][]time.Time{
			notificationType: {time.Now()},
		}
	}

	r.db[userID] = notifications

	return nil
}

func (r *usersDBRepo) GetLastNotifications(userID string, notificationType models.NotificationType) ([]time.Time, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var lastTime []time.Time

	lastNotifications, found := r.db[userID]
	if found {
		lastTime = lastNotifications[notificationType]
	}

	return lastTime, nil
}

func notificationsLessThanThreshold(lastNotifications []time.Time, threshold int) bool {
	return len(lastNotifications) >= threshold
}
