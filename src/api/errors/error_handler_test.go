package errors

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextAbortWithStatus(t *testing.T) {
	var testCases = []struct {
		description        string
		error              error
		expectedStatusCode int
	}{
		{
			description:        "Validation Error: Json is not well formed",
			error:              BadRequestError("json is not well formed"),
			expectedStatusCode: 400,
		},
		{
			description:        "Forbidden Error: notifications limit exceeded",
			error:              ForbiddenError("notifications limit exceeded"),
			expectedStatusCode: 403,
		},
		{
			description:        "Entity Not Found Error",
			error:              NotFoundError("no value was returned"),
			expectedStatusCode: 404,
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()

		t.Run(tc.description, func(t *testing.T) {
			Handle(w, tc.error)
			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
