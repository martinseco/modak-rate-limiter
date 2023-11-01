package services

import (
	"github.com/stretchr/testify/mock"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	rateLimiterMock struct {
		mock.Mock
	}

	senderMock struct {
		mock.Mock
	}
)

func (m *rateLimiterMock) Allow(userID string, notificationType models.NotificationType) (bool, error) {
	args := m.Called(userID, notificationType)

	if args.Get(1) == nil {
		return args.Get(0).(bool), nil
	}

	return false, args.Get(1).(error)
}

func (m *senderMock) Send(request models.NotificationRequest) error {
	args := m.Called(request)

	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}
