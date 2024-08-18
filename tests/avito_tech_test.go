package tests

import (
	"avito_tech/internal/entity"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

var (
	tokenModerator = ""
	tokenClient    = ""
)

func TestAvitoTechDummyLogin(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	response := e.GET("/dummyLogin").
		WithQuery("user_type", "moderator").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token").
		Value("token").String().NotEmpty()

	tokenModerator = response.Raw()

	if !regexp.MustCompile(`^[\w-]+\.[\w-]+\.[\w-]+$`).MatchString(tokenModerator) {
		t.Fatalf("token does not match the expected format: %s", tokenModerator)
	}

	response = e.GET("/dummyLogin").
		WithQuery("user_type", "client").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token").
		Value("token").String().NotEmpty()

	tokenClient = response.Raw()

	if !regexp.MustCompile(`^[\w-]+\.[\w-]+\.[\w-]+$`).MatchString(tokenClient) {
		t.Fatalf("token does not match the expected format: %s", tokenClient)
	}

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

func TestAvitoTechFlatUpdate(t *testing.T) {
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

	tokenNewModerator := response.Value("token").String().Raw()

	resp := e.POST("/flat/create").
		WithHeader("Authorization", "Bearer "+tokenModerator).
		WithJSON(entity.Flat{
			HouseID: 7,
			Number:  197,
			Rooms:   4,
			Price:   8900000,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	idFlat := resp.Value("id").Number().Raw()

	testCases := []struct {
		name    string
		status  int
		token   string
		error   string
		message string
		request entity.Flat
	}{
		{
			name:   "moderator update flat",
			status: http.StatusOK,
			token:  tokenModerator,
			request: entity.Flat{
				ID:      int64(idFlat),
				HouseID: 7,
				Number:  197,
				Rooms:   4,
				Price:   8900000,
				Status:  "on moderation",
			},
		},
		{
			name:    "client update flat",
			status:  http.StatusForbidden,
			message: "Forbidden",
			token:   tokenClient,
			request: entity.Flat{
				ID:      int64(idFlat),
				HouseID: 7,
				Number:  91,
				Rooms:   2,
				Price:   11000000,
				Status:  "on moderation",
			},
		},
		{
			name:    "new moderator update flat",
			status:  http.StatusInternalServerError,
			message: "failed to update flat",
			token:   tokenNewModerator,
			request: entity.Flat{
				ID:      int64(idFlat),
				HouseID: 7,
				Number:  91,
				Rooms:   2,
				Price:   11000000,
				Status:  "on moderation",
			},
		},
		{
			name:    "not id house",
			status:  http.StatusInternalServerError,
			token:   tokenModerator,
			message: "failed to update flat",
			request: entity.Flat{
				ID:     int64(idFlat),
				Number: 197,
				Rooms:  3,
				Price:  8900000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := httpexpect.Default(t, u.String())

			resp := e.POST("/flat/update").
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
