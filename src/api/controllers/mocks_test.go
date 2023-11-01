package controllers

import (
	"github.com/stretchr/testify/mock"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	notificationsServiceMock struct {
		mock.Mock
	}
)

func (m *notificationsServiceMock) SendNotification(request models.NotificationRequest) error {
	args := m.Called(request)

	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}
