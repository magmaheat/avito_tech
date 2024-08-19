package house_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/house"
	"avito_tech/internal/http_server/handlers/house/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateH(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		modeCreateFunc  int
		requestBody     interface{}
	}{
		{
			name:           "create House",
			expectedStatus: http.StatusOK,
			modeCreateFunc: 1,
			requestBody:    entity.House{},
		},
		{
			name:            "error decode",
			expectedMessage: "failed to add house",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   fmt.Errorf("mock error"),
			modeCreateFunc:  2,
			requestBody:     entity.House{},
		},
		{
			name:            "failed decode",
			expectedMessage: "failed to decode request body",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   fmt.Errorf("mock error"),
			requestBody:     entity.User{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewHouseStorage(t)

			switch tt.modeCreateFunc {
			case 1:
				storageMock.On("CreateH", entity.House{}).
					Return(int64(3), nil).Once()
			case 2:
				storageMock.On("CreateH", mock.Anything).
					Return(int64(-1), tt.expectedError).Once()
			}

			handler := house.Create(nil, storageMock)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/house/create", bytes.NewReader(input))
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

func TestGetAllFlats(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		id              string
		role            string
		modeCreateFunc  int
	}{
		{
			name:           "success Get Flats",
			id:             "1",
			role:           "moderator",
			expectedStatus: http.StatusOK,
			modeCreateFunc: 1,
		},
		{
			name:            "error decode request",
			id:              "1",
			role:            "moderator",
			expectedMessage: "failed to get flats",
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   fmt.Errorf("mock error"),
			modeCreateFunc:  2,
		},
		{
			name:            "empty id",
			id:              "",
			role:            "moderator",
			expectedMessage: "id is empty",
			expectedStatus:  http.StatusBadRequest,
			expectedError:   fmt.Errorf("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storageMock := mocks.NewHouseStorage(t)

			switch tt.modeCreateFunc {
			case 1:
				storageMock.On("GetAllFlats", mock.Anything, mock.Anything).
					Return([]entity.Flat{}, nil).Once()
			case 2:
				storageMock.On("GetAllFlats", mock.Anything, mock.Anything).
					Return(nil, tt.expectedError).Once()
			case 3:

			}

			handler := house.GetAllFlats(nil, storageMock)

			r := chi.NewRouter()
			r.Get("/house/{id}", handler)
			r.Get("/house/", handler)

			req, err := http.NewRequest(http.MethodGet, "/house/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), "role", tt.role)
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

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
		expectedError   error
		id              string
		modeCreateFunc  int
		requestBody     interface{}
	}{
		{
			name:           "success subscribe",
			id:             "1",
			expectedStatus: http.StatusOK,
			modeCreateFunc: 1,
			requestBody:    entity.Subscription{},
		},
		{
			name:            "empty id",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "house_id is required",
			expectedError:   fmt.Errorf("mock error"),
		},
		{
			name:            "failed decode",
			id:              "1",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "failed to decode",
			expectedError:   fmt.Errorf("mock error"),
			requestBody:     entity.Flat{},
			modeCreateFunc:  3,
		},
		{
			name:            "failed house id",
			id:              "1",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid house_id",
			expectedError:   fmt.Errorf("mock error"),
			requestBody:     entity.Subscription{},
			modeCreateFunc:  4,
		},
		{
			name:            "failed subscribe",
			id:              "1",
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "failed to subscribe",
			expectedError:   fmt.Errorf("mock error"),
			requestBody:     entity.Subscription{},
			modeCreateFunc:  2,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			storageMock := mocks.NewHouseStorage(t)

			var patches *gomonkey.Patches

			switch tt.modeCreateFunc {
			case 1:
				storageMock.On("Subscribe", mock.Anything).
					Return(nil).Once()
			case 2:
				storageMock.On("Subscribe", mock.Anything).
					Return(tt.expectedError).Once()
			case 3:
				patches = gomonkey.ApplyFunc(render.DecodeJSON, func(r io.Reader, v interface{}) error {
					return tt.expectedError
				})
				defer patches.Reset()
			case 4:
				patches = gomonkey.ApplyFunc(strconv.Atoi, func(s string) (int, error) {
					return 0, tt.expectedError
				})
				defer patches.Reset()
			}

			handler := house.Subscribe(nil, storageMock)

			r := chi.NewRouter()
			r.Post("/house/{id}/subscribe", handler)
			r.Post("/house//subscribe", handler)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/house/"+tt.id+"/subscribe", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

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
