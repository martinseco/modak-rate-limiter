package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

func TestGetLastNotifications(t *testing.T) {
	cases := []struct {
		name             string
		userID           string
		notificationType models.NotificationType
		db               map[string]models.LastNotifications
		checker          func([]time.Time, error)
	}{
		{
			name:             "when user has no notifications then return nil notifications slice",
			userID:           "not_found",
			notificationType: models.News,
			db:               map[string]models.LastNotifications{},
			checker: func(times []time.Time, err error) {
				assert.Nil(t, times)
				assert.Nil(t, err)
			},
		},
		{
			name:             "when user has notifications but not for specified type then return nil notifications slice",
			userID:           "user_123",
			notificationType: models.News,
			db: map[string]models.LastNotifications{
				"user_123": map[models.NotificationType][]time.Time{
					models.Status: {time.Now()},
				},
			},
			checker: func(times []time.Time, err error) {
				assert.Nil(t, times)
				assert.Nil(t, err)
			},
		},
		{
			name:             "when user has notifications for specified type then return notifications slice",
			userID:           "user_123",
			notificationType: models.News,
			db: map[string]models.LastNotifications{
				"user_123": map[models.NotificationType][]time.Time{
					models.News: {time.Now()},
				},
			},
			checker: func(times []time.Time, err error) {
				assert.NotNil(t, times)
				assert.Nil(t, err)
				assert.Len(t, times, 1)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewUsersDBRepo(tc.db)

			tc.checker(repo.GetLastNotifications(tc.userID, tc.notificationType))
		})
	}
}

func TestInsertNotification(t *testing.T) {
	var (
		userID          = "user_123"
		oldNotification = time.Now().AddDate(0, 0, -1)
	)

	cases := []struct {
		name             string
		userID           string
		notificationType models.NotificationType
		rule             models.NotificationRule
		db               map[string]models.LastNotifications
		resultChecker    func(error)
		dbChecker        func(map[string]models.LastNotifications, map[string]models.LastNotifications)
	}{
		{
			name:             "when it is first notification for user then insert notification",
			userID:           userID,
			notificationType: models.News,
			rule:             models.NotificationRule{AmountThreshold: 1},
			db:               map[string]models.LastNotifications{},
			resultChecker: func(err error) {
				assert.Nil(t, err)
			},
			dbChecker: func(initialDB map[string]models.LastNotifications, resultDB map[string]models.LastNotifications) {
				assert.NotEqual(t, initialDB, resultDB)
			},
		},
		{
			name:             "when it is first notification for type then insert notification",
			userID:           userID,
			notificationType: models.News,
			rule:             models.NotificationRule{AmountThreshold: 1},
			db: map[string]models.LastNotifications{
				"user_123": map[models.NotificationType][]time.Time{
					models.Marketing: {time.Now()},
				},
			},
			resultChecker: func(err error) {
				assert.Nil(t, err)
			},
			dbChecker: func(initialDB map[string]models.LastNotifications, resultDB map[string]models.LastNotifications) {
				assert.NotEqual(t, initialDB, resultDB)
			},
		},
		{
			name:             "when notifications are less than amount threshold for type then insert notification",
			userID:           userID,
			notificationType: models.News,
			rule:             models.NotificationRule{AmountThreshold: 2},
			db: map[string]models.LastNotifications{
				"user_123": map[models.NotificationType][]time.Time{
					models.News: {time.Now()},
				},
			},
			resultChecker: func(err error) {
				assert.Nil(t, err)
			},
			dbChecker: func(initialDB map[string]models.LastNotifications, resultDB map[string]models.LastNotifications) {
				assert.NotEqual(t, initialDB, resultDB)
				assert.Equal(t, initialDB[userID][models.News][0], resultDB[userID][models.News][0])
			},
		},
		{
			name:             "when notifications quantity is equal than amount threshold for type then remove oldest and insert new notification",
			userID:           userID,
			notificationType: models.News,
			rule:             models.NotificationRule{AmountThreshold: 2},
			db: map[string]models.LastNotifications{
				"user_123": map[models.NotificationType][]time.Time{
					models.News: {oldNotification, time.Now()},
				},
			},
			resultChecker: func(err error) {
				assert.Nil(t, err)
			},
			dbChecker: func(initialDB map[string]models.LastNotifications, resultDB map[string]models.LastNotifications) {
				assert.NotEqual(t, initialDB, resultDB)
				assert.NotEqual(t, initialDB[userID][models.News][0], resultDB[userID][models.News][0])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewUsersDBRepo(tc.db)

			initialDBCopy := copyDB(tc.db)

			tc.resultChecker(repo.InsertNotification(tc.userID, tc.notificationType, tc.rule))

			tc.dbChecker(initialDBCopy, tc.db)
		})
	}
}

func copyDB(db map[string]models.LastNotifications) map[string]models.LastNotifications {
	initialDBCopy := map[string]models.LastNotifications{}
	notificationsCopy := map[models.NotificationType][]time.Time{}

	for user, notifications := range db {
		for notificationType, times := range notifications {
			notificationsCopy[notificationType] = times
		}
		initialDBCopy[user] = notificationsCopy
	}
	return initialDBCopy
}
