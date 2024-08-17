package tests

import (
	"avito_tech/internal/entity"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

func TestAvitoTechDummyLogin(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.GET("/dummyLogin").
		WithQuery("user_type", "moderator").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token").
		Value("token").String().NotEmpty().
		Match(`^[\w-]+\.[\w-]+\.[\w-]+$`)

	e.GET("/dummyLogin").
		WithQuery("user_type", "client").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token").
		Value("token").String().NotEmpty().
		Match(`^[\w-]+\.[\w-]+\.[\w-]+$`)

	e.GET("/dummyLogin").
		WithQuery("user_type", "hacker").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()
}

func TestAvitoTechHouseCreate(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	response := e.GET("/dummyLogin").
		WithQuery("user_type", "moderator").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tokenModerator := response.Value("token").String().NotEmpty().Raw()

	response = e.GET("/dummyLogin").
		WithQuery("user_type", "client").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tokenClient := response.Value("token").String().NotEmpty().Raw()

	testCases := []struct {
		name    string
		status  int
		token   string
		error   string
		message string
		request entity.House
	}{
		{
			name:    "moderator created house",
			status:  http.StatusOK,
			token:   tokenModerator,
			message: "house added",
			request: entity.House{
				ID:      7,
				Address: "Moscow street, 4",
				Year:    2000,
			},
		},
		{
			name:    "not id house",
			status:  http.StatusInternalServerError,
			token:   tokenModerator,
			message: "failed to add house",
			request: entity.House{
				Address: "Moscow street, 4",
				Year:    2000,
			},
		},
		{
			name:    "client created house",
			status:  http.StatusForbidden,
			token:   tokenClient,
			message: "Forbidden",
			request: entity.House{
				ID:      143,
				Address: "Avito street, 143",
				Year:    1997,
			},
		},
		{
			name:    "not token",
			status:  http.StatusForbidden,
			message: "failed check sing token",
			token:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := httpexpect.Default(t, u.String())

			e.POST("/house/create").
				WithHeader("Authorization", "Bearer "+tc.token).
				WithJSON(tc.request).
				Expect().
				Status(tc.status).
				JSON().Object().
				ContainsKey("message").
				Value("message").String().Match(tc.message)
		})
	}
}

func TestAvitoTechFlatsCreate(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	response := e.GET("/dummyLogin").
		WithQuery("user_type", "moderator").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tokenModerator := response.Value("token").String().NotEmpty().Raw()

	response = e.GET("/dummyLogin").
		WithQuery("user_type", "client").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tokenClient := response.Value("token").String().NotEmpty().Raw()

	testCases := []struct {
		name    string
		status  int
		token   string
		error   string
		message string
		request entity.Flat
	}{
		{
			name:   "moderator create flat",
			status: http.StatusOK,
			token:  tokenModerator,
			request: entity.Flat{
				HouseID: 7,
				Number:  197,
				Rooms:   3,
				Price:   8900000,
			},
		},
		{
			name:   "client create flat",
			status: http.StatusOK,
			token:  tokenClient,
			request: entity.Flat{
				HouseID: 7,
				Number:  91,
				Rooms:   2,
				Price:   11000000,
			},
		},
		{
			name:    "not id house",
			status:  http.StatusInternalServerError,
			token:   tokenModerator,
			message: "failed to add flat",
			request: entity.Flat{
				HouseID: 99,
				Number:  197,
				Rooms:   3,
				Price:   8900000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := httpexpect.Default(t, u.String())

			resp := e.POST("/flat/create").
				WithHeader("Authorization", "Bearer "+tc.token).
				WithJSON(tc.request).
				Expect().
				Status(tc.status).
				JSON().Object()

			if tc.message != "" {
				resp.Value("message").String().IsEqual(tc.message)
			}
		})
	}
}
