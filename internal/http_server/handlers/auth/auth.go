package auth

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/lib/logger/slg"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Storage interface {
	Check(email entity.User) bool
	Register()
}

var MySigningKey = []byte(os.Getenv("MY_SIGNING_KEY"))

func DummyLogin(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.DummyLogin"
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

		if user.UserType != "moderator" {
			user.UserType = "client"
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Email,
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

		render.JSON(w, r, map[string]string{"token": tokenString})
	}
}

func Register(log *slog.Logger, storage Storage) http.HandlerFunc {
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

func JWTAuth(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.JWTModerator"
		reqID := middleware.GetReqID(r.Context())

		log = slg.SetupLogger(fn, reqID)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			message := "failed check sing token"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": message})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role, ok := claims["role"].(string)

			if !ok {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"error": "Forbidden"})
				return
			}

			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Error("message", slg.Err(err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized"})
		}
	}
}

func RequireModerator(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.JWTRole"
		reqID := middleware.GetReqID(r.Context())
		role := r.Context().Value("role").(string)

		log = slg.SetupLogger(fn, reqID)

		if role != "moderator" {
			message := "Forbidden"
			log.Error(message)
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": message})
			return
		}

		next.ServeHTTP(w, r)
	}
}
