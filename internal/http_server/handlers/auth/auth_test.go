package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/auth"
	"avito_tech/internal/http_server/handlers/auth/mocks"
)

func TestAuthHandler(t *testing.T) {
	cases := []struct {
		name            string
		userType        string
		mockError       error
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:           "Success Client",
			userType:       "client",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Success Moderator",
			userType:       "moderator",
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Error Creating User",
			userType:        "client",
			expectedMessage: "failed added user",
		},
		{
			name:            "Error Generating Hash Password",
			userType:        "client",
			expectedMessage: "failed to generate hash password",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewAuthStorage(t)

			if tc.mockError != nil {
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.Nil, tc.mockError).Once()
			} else {
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.New(), nil).Once()
			}

			handler := auth.DummyLogin(nil, storageMock)

			user := entity.User{
				UserType: tc.userType,
			}

			input, err := json.Marshal(user)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/dummyLogin", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedMessage != "" {
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, tc.expectedMessage, response["message"])
			}
		})
	}
}
