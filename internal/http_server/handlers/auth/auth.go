package auth

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/lib/logger/slg"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=AuthStorage
type AuthStorage interface {
	CreateUser(user entity.User, hashPassword string) (uuid.UUID, error)
}

type ResponseDummyLogin struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	UserType string    `json:"user_type"`
	Token    string    `json:"token"`
}

var MySigningKey = []byte(os.Getenv("MY_SIGNING_KEY"))

func DummyLogin(log *slog.Logger, storage AuthStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.DummyLogin"
		reqID := middleware.GetReqID(r.Context())

		log = slg.SetupLogger(fn, reqID)

		userType := r.URL.Query().Get("user_type")
		if userType == "" {
			log.Error("user_type parameter is missing")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": "user_type parameter is required"})
			return
		}

		user := entity.User{
			Email:    "dummylogin@gmail.com",
			Password: "password",
			UserType: userType,
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			message := "failed to generate hash password"
			log.Error("message", message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		id, err := storage.CreateUser(user, string(hashPassword))
		if err != nil {
			message := "failed added user"
			log.Error("message", message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": id.String(),
			"role":     user.UserType,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString(MySigningKey)
		if err != nil {
			message := "failed to signed token"
			log.Error(message)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log.Info("User added", slog.Any("request", reqID))

		render.JSON(w, r, ResponseDummyLogin{
			ID:       id,
			Email:    user.Email,
			Password: user.Password,
			UserType: user.UserType,
			Token:    tokenString,
		})
	}
}

func Register(log *slog.Logger, storage AuthStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.register"
		reqID := middleware.GetReqID(r.Context())

		log = slg.SetupLogger(fn, reqID)

		var user entity.User

		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log.Info("request body decoded")

	}
}
