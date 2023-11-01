package limiter

import (
	"github.com/stretchr/testify/mock"
	"time"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	rulesDBMock struct {
		mock.Mock
	}

	usersDBMock struct {
		mock.Mock
	}
)

func (m *rulesDBMock) Insert(rule models.NotificationRule) error {
	args := m.Called(rule)

	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}

func (m *rulesDBMock) Get(notificationType models.NotificationType) (models.NotificationRule, error) {
	args := m.Called(notificationType)

	if args.Get(0) != nil {
		return args.Get(0).(models.NotificationRule), nil
	}

	return models.NotificationRule{}, args.Get(1).(error)
}

func (m *usersDBMock) InsertNotification(userID string, notificationType models.NotificationType, rule models.NotificationRule) error {
	args := m.Called(userID, notificationType, rule)

	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}

func (m *usersDBMock) GetLastNotifications(userID string, notificationType models.NotificationType) ([]time.Time, error) {
	args := m.Called(userID, notificationType)

	if args.Get(0) != nil {
		return args.Get(0).([]time.Time), nil
	}

	return nil, args.Get(1).(error)
}
