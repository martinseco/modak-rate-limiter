package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/martinseco/modak-rate-limiter/src/api/errors"
)

func TestSendNotification(t *testing.T) {
	cases := []struct {
		name             string
		reqBody          interface{}
		sendNotification func(*mock.Mock) mock.Call
		checker          func(status int, response []byte)
	}{
		{
			name:    "when request body is invalid then return status 400 and error",
			reqBody: "not_valid",
			sendNotification: func(m *mock.Mock) mock.Call {
				return *m.On("SendNotification", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(status int, response []byte) {
				assert.Equal(t, http.StatusBadRequest, status)
				err := getBodyAsError(response)
				assert.Equal(t, "json: cannot unmarshal string into Go value of type models.NotificationRequest", err.Message)
			},
		},
		{
			name: "when request is valid then return status 202",
			reqBody: models.NotificationRequest{
				UserID:  "user_123",
				Type:    models.News,
				Message: "message to send",
			},
			sendNotification: func(m *mock.Mock) mock.Call {
				return *m.On("SendNotification", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			checker: func(status int, response []byte) {
				assert.Equal(t, http.StatusAccepted, status)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := new(notificationsServiceMock)
			tc.sendNotification(&serviceMock.Mock)

			controller := NewNotificationsController(serviceMock)

			jsonBody, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(jsonBody))

			w := httptest.NewRecorder()

			controller.SendNotification(w, req)

			tc.checker(w.Code, getResponse(w))
		})
	}
}

func getResponse(recorder *httptest.ResponseRecorder) []byte {
	responseBody, _ := ioutil.ReadAll(recorder.Body)
	return responseBody
}

func getBodyAsError(data []byte) (error errors.ApiError) {
	_ = json.Unmarshal(data, &error)
	return
}
