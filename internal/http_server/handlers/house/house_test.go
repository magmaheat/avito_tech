package house_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/house"
	"avito_tech/internal/http_server/handlers/house/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateH(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		userID          uuid.UUID
	}{
		{
			name:           "Create House",
			expectedStatus: http.StatusOK,
			userID:         uuid.New(),
		},
		{
			name:            "Error decode",
			expectedMessage: "failed to add house",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   errors.New("mock error"),
			userID:          uuid.New(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewHouseStorage(t)

			if tt.expectedError != nil {
				storageMock.On("CreateH", mock.Anything).
					Return(int64(-1), tt.expectedError).Once()
			} else {
				storageMock.On("CreateH", entity.House{}).
					Return(int64(3), nil).Once()
			}

			handler := house.Create(nil, storageMock)

			houseRequest := entity.House{}

			input, err := json.Marshal(houseRequest)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(input))
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

func TestGetFlats(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		id              string
		role            string
	}{
		{
			name:           "Success Get Flats",
			id:             "1",
			role:           "moderator",
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Error decode request",
			id:              "1",
			role:            "moderator",
			expectedMessage: "failed to get flats",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewHouseStorage(t)

			if tt.expectedError != nil {
				storageMock.On("GetFlats", mock.Anything, mock.Anything).
					Return(nil, tt.expectedError).Once()
			} else {
				storageMock.On("GetFlats", mock.Anything, mock.Anything).
					Return([]entity.Flat{}, nil).Once()
			}

			handler := house.Flats(nil, storageMock)

			r := chi.NewRouter()
			r.Get("/house/{id}", handler)

			req, err := http.NewRequest(http.MethodGet, "/house/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), "role", "moderator")
			req = req.WithContext(ctx)

			r.ServeHTTP(rr, req)

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
