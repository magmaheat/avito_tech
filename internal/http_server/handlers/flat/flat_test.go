package flat_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/flat"
	"avito_tech/internal/http_server/handlers/flat/mocks"
	"avito_tech/internal/http_server/sender"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		userID          uuid.UUID
		requestBody     interface{}
		modeCreateFunc  int
	}{
		{
			name:           "Create flat",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
			requestBody:    entity.Flat{},
			modeCreateFunc: 1,
		},
		{
			name:            "failed decode",
			expectedMessage: "failed to decode request body",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   fmt.Errorf("mock error"),
			userID:          uuid.New(),
			requestBody:     entity.User{},
		},
		{
			name:           "failed get subscribers",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
			requestBody:    entity.Flat{},
			modeCreateFunc: 2,
			expectedError:  fmt.Errorf("mock error"),
		},
		{
			name:           "failed send",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
			requestBody:    entity.Flat{},
			modeCreateFunc: 3,
		},
		{
			name:            "error decode",
			expectedMessage: "failed to add flat",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   fmt.Errorf("mock error"),
			userID:          uuid.New(),
			requestBody:     entity.Flat{},
			modeCreateFunc:  4,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			storageMock := mocks.NewFlatStorage(t)

			var patches *gomonkey.Patches

			switch tt.modeCreateFunc {
			case 1:
				storageMock.On("CreateF", entity.Flat{UserID: tt.userID}).
					Return(int64(3), nil).Once()

				storageMock.On("GetSubscribers", mock.Anything).
					Return([]string{"subscriber1@example.com", "subscriber2@example.com"}, nil).Once()

			case 2:
				storageMock.On("CreateF", entity.Flat{UserID: tt.userID}).
					Return(int64(3), nil).Once()

				storageMock.On("GetSubscribers", mock.Anything).
					Return(nil, tt.expectedError).Once()

			case 3:
				storageMock.On("CreateF", entity.Flat{UserID: tt.userID}).
					Return(int64(3), nil).Once()

				storageMock.On("GetSubscribers", mock.Anything).
					Return([]string{"subscriber1@example.com", "subscriber2@example.com"}, nil).Once()

				patches = gomonkey.ApplyMethod(reflect.TypeOf((*sender.Sender)(nil)), "SendEmail", func(s *sender.Sender, ctx context.Context, email, message string) error {
					return errors.New("mock error")
				})
				defer patches.Reset()

			case 4:
				storageMock.On("CreateF", mock.Anything).
					Return(int64(-1), tt.expectedError).Once()
			}

			handler := flat.Create(nil, storageMock, nil)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), "username", tt.userID)
			req = req.WithContext(ctx)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			time.Sleep(1 * time.Second)

			if tt.expectedMessage != "" {
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, tt.expectedMessage, response["message"])
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		userID          uuid.UUID
		requestBody     interface{}
		modeCreateFunc  int
	}{
		{
			name:           "Update flat",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
			requestBody:    entity.Flat{},
			modeCreateFunc: 1,
		},
		{
			name:            "Error update",
			expectedMessage: "failed to update flat",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   errors.New("mock error"),
			userID:          uuid.New(),
			requestBody:     entity.Flat{},
			modeCreateFunc:  2,
		},
		{
			name:           "failed decode",
			expectedError:  fmt.Errorf("mock error"),
			expectedStatus: http.StatusBadRequest,
			userID:         uuid.New(),
			requestBody:    entity.User{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			storageMock := mocks.NewFlatStorage(t)

			switch tt.modeCreateFunc {
			case 1:
				storageMock.On("Update", mock.Anything, tt.userID).
					Return(nil).Once()
			case 2:
				storageMock.On("Update", mock.Anything, mock.Anything).
					Return(tt.expectedError).Once()
			}

			handler := flat.Update(nil, storageMock)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/flat/update", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), "username", tt.userID)
			req = req.WithContext(ctx)

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
