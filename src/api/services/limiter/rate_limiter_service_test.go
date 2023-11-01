package limiter

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	customErrors "github.com/martinseco/modak-rate-limiter/src/api/errors"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
	"github.com/stretchr/testify/mock"
)

func TestAllow(t *testing.T) {
	cases := []struct {
		name                string
		userID              string
		notificationType    models.NotificationType
		getRules            func(mock *mock.Mock) mock.Call
		getLastNotification func(mock *mock.Mock) mock.Call
		insertNotification  func(mock *mock.Mock) mock.Call
		checker             func(bool, error)
	}{
		{
			name:             "when rule not found for notification type then return error",
			userID:           "user_123",
			notificationType: "not_found",
			getRules: func(m *mock.Mock) mock.Call {
				return *m.On("Get", mock.Anything, mock.Anything).
					Return(nil, customErrors.NotFoundError("")).Once()
			},
			getLastNotification: func(m *mock.Mock) mock.Call {
				return *m.On("GetLastNotifications", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil).Once()
			},
			insertNotification: func(m *mock.Mock) mock.Call {
				return *m.On("InsertNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(response bool, err error) {
				assert.False(t, response)
				assert.NotNil(t, err)
				assert.Equal(t, reflect.TypeOf(customErrors.NotFoundError("")), reflect.TypeOf(err))
			},
		},
		{
			name:             "when is first notification for userID and notificationType then return true",
			userID:           "user_123",
			notificationType: models.News,
			getRules: func(m *mock.Mock) mock.Call {
				return *m.On("Get", mock.Anything, mock.Anything).
					Return(
						models.NotificationRule{
							AmountThreshold:   1,
							DurationThreshold: 10 * time.Second,
						}, nil).Once()
			},
			getLastNotification: func(m *mock.Mock) mock.Call {
				return *m.On("GetLastNotifications", mock.Anything, mock.Anything, mock.Anything).
					Return([]time.Time{}, nil).Once()
			},
			insertNotification: func(m *mock.Mock) mock.Call {
				return *m.On("InsertNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(response bool, err error) {
				assert.True(t, response)
				assert.Nil(t, err)
			},
		},
		{
			name:             "when last notification exceeds threshold then return true",
			userID:           "user_123",
			notificationType: models.News,
			getRules: func(m *mock.Mock) mock.Call {
				return *m.On("Get", mock.Anything, mock.Anything).
					Return(
						models.NotificationRule{
							AmountThreshold:   1,
							DurationThreshold: 1 * time.Nanosecond,
						}, nil).Once()
			},
			getLastNotification: func(m *mock.Mock) mock.Call {
				return *m.On("GetLastNotifications", mock.Anything, mock.Anything, mock.Anything).
					Return([]time.Time{time.Now()}, nil).Once()
			},
			insertNotification: func(m *mock.Mock) mock.Call {
				return *m.On("InsertNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(response bool, err error) {
				assert.True(t, response)
				assert.Nil(t, err)
			},
		},
		{
			name:             "when last notification is within threshold then return false",
			userID:           "user_123",
			notificationType: models.News,
			getRules: func(m *mock.Mock) mock.Call {
				return *m.On("Get", mock.Anything, mock.Anything).
					Return(
						models.NotificationRule{
							AmountThreshold:   1,
							DurationThreshold: 10 * time.Second,
						}, nil).Once()
			},
			getLastNotification: func(m *mock.Mock) mock.Call {
				return *m.On("GetLastNotifications", mock.Anything, mock.Anything, mock.Anything).
					Return([]time.Time{time.Now()}, nil).Once()
			},
			insertNotification: func(m *mock.Mock) mock.Call {
				return *m.On("InsertNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(response bool, err error) {
				assert.False(t, response)
				assert.Nil(t, err)
			},
		},
		{
			name:             "when inserting last notification fails then return error",
			userID:           "user_123",
			notificationType: models.News,
			getRules: func(m *mock.Mock) mock.Call {
				return *m.On("Get", mock.Anything, mock.Anything).
					Return(
						models.NotificationRule{
							AmountThreshold:   1,
							DurationThreshold: 1 * time.Nanosecond,
						}, nil).Once()
			},
			getLastNotification: func(m *mock.Mock) mock.Call {
				return *m.On("GetLastNotifications", mock.Anything, mock.Anything, mock.Anything).
					Return([]time.Time{time.Now()}, nil).Once()
			},
			insertNotification: func(m *mock.Mock) mock.Call {
				return *m.On("InsertNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("error inserting record")).Once()
			},
			checker: func(response bool, err error) {
				assert.False(t, response)
				assert.NotNil(t, err)
				assert.Equal(t, "error inserting record", err.Error())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rulesDBService := new(rulesDBMock)
			usersDBService := new(usersDBMock)

			tc.getRules(&rulesDBService.Mock)
			tc.getLastNotification(&usersDBService.Mock)
			tc.insertNotification(&usersDBService.Mock)

			service := NewRateLimiterService(rulesDBService, usersDBService)

			tc.checker(service.Allow(tc.userID, tc.notificationType))
		})
	}
}

func TestIsFirstNotification(t *testing.T) {
	cases := []struct {
		name          string
		notifications []time.Time
		checker       func(result bool)
	}{
		{
			name:          "when notifications is nil then return true",
			notifications: nil,
			checker: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name:          "when notifications is empty then return true",
			notifications: []time.Time{},
			checker: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name:          "when notifications is not empty then return true",
			notifications: []time.Time{time.Now()},
			checker: func(result bool) {
				assert.False(t, result)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checker(isFirstNotification(tc.notifications))
		})
	}
}

func TestIsOverLimit(t *testing.T) {
	cases := []struct {
		name          string
		rule          models.NotificationRule
		notifications []time.Time
		checker       func(result bool)
	}{
		{
			name:          "when there are less notifications than rule amount threshold then return false",
			rule:          models.NotificationRule{AmountThreshold: 3},
			notifications: []time.Time{time.Now()},
			checker: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "when oldest notification is within the duration threshold then return true",
			rule: models.NotificationRule{
				AmountThreshold:   1,
				DurationThreshold: 1 * time.Hour,
			},
			notifications: []time.Time{time.Now()},
			checker: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "when oldest notification exceeds the duration threshold then return false",
			rule: models.NotificationRule{
				AmountThreshold:   1,
				DurationThreshold: 1 * time.Nanosecond,
			},
			notifications: []time.Time{time.Now()},
			checker: func(result bool) {
				assert.False(t, result)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checker(isOverLimit(tc.notifications, tc.rule))
		})
	}
}
