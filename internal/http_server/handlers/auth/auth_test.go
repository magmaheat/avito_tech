package auth_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/auth"
	"avito_tech/internal/http_server/handlers/auth/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDummyLogin(t *testing.T) {
	tests := []struct {
		name            string
		userType        string
		expectedMessage string
		expectedStatus  int
		mockError       error
	}{
		{
			name:           "Success Client",
			userType:       "client",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SuccessModerator",
			userType:       "moderator",
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Error Creating User",
			userType:        "client",
			expectedMessage: "failed added user",
			expectedStatus:  http.StatusInternalServerError,
			mockError:       errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewAuthStorage(t)

			if tt.mockError != nil {
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.Nil, tt.mockError).Once()
			} else {
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.New(), nil).Once()
			}

			handler := auth.DummyLogin(nil, storageMock)

			user := entity.User{
				UserType: tt.userType,
			}

			input, err := json.Marshal(user)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/dummyLogin", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedMessage != "" {
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, tt.expectedMessage, response["message"])
			}
		})
	}
}

//func TestRegister(t *testing.T) {
//	type args struct {
//		log     *slog.Logger
//		storage AuthStorage
//	}
//	tests := []struct {
//		name string
//		args args
//		want http.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := Register(tt.args.log, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Register() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
