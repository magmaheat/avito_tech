package auth_test

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/handlers/auth"
	"avito_tech/internal/http_server/handlers/auth/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDummyLogin(t *testing.T) {
	tests := []struct {
		name               string
		userType           string
		expectedMessage    string
		expectedStatus     int
		modeCreateMockFunc int
		mockError          error
	}{
		{
			name:               "Success Client",
			userType:           "client",
			expectedStatus:     http.StatusOK,
			modeCreateMockFunc: 1,
		},
		{
			name:               "SuccessModerator",
			userType:           "moderator",
			expectedStatus:     http.StatusOK,
			modeCreateMockFunc: 1,
		},
		{
			name:               "Error Creating User",
			userType:           "hacker",
			expectedMessage:    "failed added user",
			expectedStatus:     http.StatusInternalServerError,
			modeCreateMockFunc: -1,
			mockError:          fmt.Errorf("mock error"),
		},
		{
			name:               "No user_type",
			expectedMessage:    "user_type parameter is required",
			expectedStatus:     http.StatusBadRequest,
			modeCreateMockFunc: 0,
		},
		{
			name:               "hash password",
			userType:           "moderator",
			expectedMessage:    "failed to generate hash password",
			expectedStatus:     http.StatusInternalServerError,
			modeCreateMockFunc: 2,
			mockError:          fmt.Errorf("mock error"),
		},
		{
			name:               "sign key",
			userType:           "moderator",
			expectedMessage:    "failed to signed token",
			expectedStatus:     http.StatusInternalServerError,
			modeCreateMockFunc: 3,
			mockError:          fmt.Errorf("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewAuthStorage(t)
			var patches *gomonkey.Patches

			switch tt.modeCreateMockFunc {
			case 1:
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.New(), nil).Once()
			case -1:
				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.Nil, tt.mockError).Once()
			case 2:
				patches = gomonkey.ApplyFunc(bcrypt.GenerateFromPassword, func(password []byte, cost int) ([]byte, error) {
					return nil, tt.mockError
				})
				defer patches.Reset()
			case 3:
				patches = gomonkey.ApplyFunc(jwt.NewWithClaims, func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
					return &jwt.Token{
						Method: method,
						Claims: claims,
					}
				})
				defer patches.Reset()

				patches = gomonkey.ApplyMethod(reflect.TypeOf((*jwt.Token)(nil)), "SignedString", func(token *jwt.Token, key interface{}) (string, error) {
					return "", errors.New("mock error")
				})
				defer patches.Reset()

				storageMock.On("CreateUser", mock.Anything, mock.Anything).
					Return(uuid.New(), nil).Once()
			}

			handler := auth.DummyLogin(nil, storageMock)

			user := entity.User{
				UserType: tt.userType,
			}

			input, err := json.Marshal(user)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/dummyLogin?user_type="+user.UserType, bytes.NewReader(input))
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

func TestRegister(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatus     int
		expectedMessage    string
		modeCreateMockFunc int
		mockError          error
		requestBody        interface{}
	}{
		{
			name:               "register user",
			expectedStatus:     http.StatusOK,
			modeCreateMockFunc: 1,
			requestBody:        entity.User{},
		},
		{
			name:            "failed to decode",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "failed to decode request body",
			requestBody:     entity.House{},
		},
		{
			name:               "generate hash password",
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "failed to generate hash password",
			modeCreateMockFunc: 2,
			requestBody:        entity.User{},
			mockError:          fmt.Errorf("mock error"),
		},
		{
			name:               "failed register",
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "failed to register user",
			modeCreateMockFunc: -1,
			requestBody:        entity.User{},
			mockError:          fmt.Errorf("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewAuthStorage(t)
			var patches *gomonkey.Patches

			switch tt.modeCreateMockFunc {
			case 1:
				storageMock.On("Register", mock.Anything).
					Return(uuid.New().String(), nil).Once()
			case 2:
				patches = gomonkey.ApplyFunc(bcrypt.GenerateFromPassword, func(password []byte, cost int) ([]byte, error) {
					return nil, tt.mockError
				})
				defer patches.Reset()
			case -1:
				storageMock.On("Register", mock.Anything).
					Return("", tt.mockError).Once()
			}

			handler := auth.Register(nil, storageMock)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(input))
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

func TestLogin(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatus     int
		expectedMessage    string
		modeCreateMockFunc int
		mockError          error
		requestBody        interface{}
	}{
		{
			name:               "success login user",
			expectedStatus:     http.StatusOK,
			modeCreateMockFunc: 1,
			requestBody:        entity.User{},
		},
		{
			name:            "failed decode",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "failed to decode request body",
			requestBody:     entity.House{},
		},
		{
			name:               "not found",
			expectedStatus:     http.StatusNotFound,
			expectedMessage:    "user not found",
			requestBody:        entity.User{},
			modeCreateMockFunc: 3,
			mockError:          errors.New("user not found: storage.postgres.Login"),
		},
		{
			name:               "failed login",
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "failed to build query",
			requestBody:        entity.User{},
			modeCreateMockFunc: 3,
			mockError:          errors.New("mock error"),
		},
		{
			name:               "unauthorized",
			expectedStatus:     http.StatusUnauthorized,
			expectedMessage:    "invalid password",
			requestBody:        entity.User{},
			modeCreateMockFunc: 4,
			mockError:          fmt.Errorf("mock error"),
		},
		{
			name:               "bad token",
			expectedStatus:     http.StatusInternalServerError,
			expectedMessage:    "failed to signed token",
			requestBody:        entity.User{},
			modeCreateMockFunc: 2,
			mockError:          fmt.Errorf("mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewAuthStorage(t)
			var patches *gomonkey.Patches

			_ = patches

			switch tt.modeCreateMockFunc {
			case 1:
				storageMock.On("Login", mock.Anything).
					Return(entity.User{}, nil).Once()

				patches = gomonkey.ApplyFunc(bcrypt.CompareHashAndPassword, func(storagePassword []byte, password []byte) error {
					return nil
				})
				defer patches.Reset()

			case 2:
				storageMock.On("Login", mock.Anything).
					Return(entity.User{}, nil).Once()

				patches = gomonkey.ApplyFunc(bcrypt.CompareHashAndPassword, func(storagePassword []byte, password []byte) error {
					return nil
				})
				defer patches.Reset()

				patches = gomonkey.ApplyFunc(jwt.NewWithClaims, func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
					return &jwt.Token{
						Method: method,
						Claims: claims,
					}
				})
				defer patches.Reset()

				patches = gomonkey.ApplyMethod(reflect.TypeOf((*jwt.Token)(nil)), "SignedString", func(token *jwt.Token, key interface{}) (string, error) {
					return "", errors.New("mock error")
				})
				defer patches.Reset()

			case 3:
				storageMock.On("Login", mock.Anything).
					Return(entity.User{}, tt.mockError).Once()

			case 4:
				storageMock.On("Login", mock.Anything).
					Return(entity.User{}, nil).Once()
			}

			handler := auth.Login(nil, storageMock)

			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(input))
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
