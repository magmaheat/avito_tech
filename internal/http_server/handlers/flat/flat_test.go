package flat_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/flat"
	"avito_tech/internal/http_server/handlers/flat/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		userID          uuid.UUID
	}{
		{
			name:           "Create flat",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
		},
		{
			name:            "error decode",
			expectedMessage: "failed to add flat",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   errors.New("mock error"),
			userID:          uuid.New(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewFlatStorage(t)

			if tt.expectedError != nil {
				storageMock.On("CreateF", mock.Anything).
					Return(int64(-1), tt.expectedError).Once()
			} else {
				storageMock.On("CreateF", entity.Flat{UserID: tt.userID}).
					Return(int64(3), nil).Once()
			}

			handler := flat.Create(nil, storageMock)

			flatRequest := entity.Flat{
				UserID: tt.userID,
			}

			input, err := json.Marshal(flatRequest)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/flat/create", bytes.NewReader(input))
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

func TestUpdate(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		userID          uuid.UUID
	}{
		{
			name:           "Update flat",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
		},
		{
			name:            "Error update",
			expectedMessage: "failed to update flat",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   errors.New("mock error"),
			userID:          uuid.New(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewFlatStorage(t)

			if tt.expectedError != nil {
				storageMock.On("Update", mock.Anything, mock.Anything).
					Return(tt.expectedError).Once()
			} else {
				storageMock.On("Update", entity.Flat{UserID: tt.userID}, tt.userID).
					Return(nil).Once()
			}

			handler := flat.Update(nil, storageMock)

			flatRequest := entity.Flat{
				UserID: tt.userID,
			}

			input, err := json.Marshal(flatRequest)
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
