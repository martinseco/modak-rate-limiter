package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	internalErrors "github.com/martinseco/modak-rate-limiter/src/api/errors"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

func TestSendNotification(t *testing.T) {
	cases := []struct {
		name    string
		request models.NotificationRequest
		allow   func(*mock.Mock) mock.Call
		send    func(*mock.Mock) mock.Call
		checker func(error)
	}{
		{
			name:    "when rate limiter returns error then return error",
			request: models.NotificationRequest{},
			allow: func(m *mock.Mock) mock.Call {
				return *m.On("Allow", mock.Anything, mock.Anything, mock.Anything).
					Return(false, errors.New("unexpected error on rateLimiter.Allow")).Once()
			},
			send: func(m *mock.Mock) mock.Call {
				return *m.On("Send", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "unexpected error on rateLimiter.Allow", err.Error())
			},
		},
		{
			name:    "when allowed is false then return forbidden error",
			request: models.NotificationRequest{UserID: "userID", Type: models.News},
			allow: func(m *mock.Mock) mock.Call {
				return *m.On("Allow", mock.Anything, mock.Anything, mock.Anything).
					Return(false, nil).Once()
			},
			send: func(m *mock.Mock) mock.Call {
				return *m.On("Send", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, reflect.TypeOf(internalErrors.ForbiddenError("")), reflect.TypeOf(err))
			},
		},
		{
			name:    "when allowed is true then return send notification",
			request: models.NotificationRequest{UserID: "userID", Type: models.News},
			allow: func(m *mock.Mock) mock.Call {
				return *m.On("Allow", mock.Anything, mock.Anything, mock.Anything).
					Return(true, nil).Once()
			},
			send: func(m *mock.Mock) mock.Call {
				return *m.On("Send", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name:    "when send notification fails then return error",
			request: models.NotificationRequest{UserID: "userID", Type: models.News},
			allow: func(m *mock.Mock) mock.Call {
				return *m.On("Allow", mock.Anything, mock.Anything, mock.Anything).
					Return(true, nil).Once()
			},
			send: func(m *mock.Mock) mock.Call {
				return *m.On("Send", mock.Anything, mock.Anything).
					Return(errors.New("unexpected error on sender.Send")).Once()
			},
			checker: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "unexpected error on sender.Send", err.Error())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rateLimiterService := new(rateLimiterMock)
			senderService := new(senderMock)

			tc.allow(&rateLimiterService.Mock)
			tc.send(&senderService.Mock)

			service := NewNotificationsService(rateLimiterService, senderService)

			tc.checker(service.SendNotification(tc.request))
		})
	}
}
