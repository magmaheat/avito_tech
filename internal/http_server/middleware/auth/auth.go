package auth

import (
	"avito_tech/internal/lib/logger/slg"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var MySigningKey = []byte(os.Getenv("MY_SIGNING_KEY"))

func JWTAuth(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.JWTModerator"
		reqID := middleware.GetReqID(r.Context())

		log = slg.WithLogger(fn, reqID)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error("Unauthorized")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"message": "Unauthorized"})
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
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role, okRole := claims["role"].(string)
			usernameString, okName := claims["username"].(string)

			if !okRole || !okName {
				log.Error("Forbidden")
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"message": "Forbidden"})
				return
			}

			username, _ := uuid.Parse(usernameString)

			ctx := context.WithValue(r.Context(), "role", role)
			ctx = context.WithValue(ctx, "username", username)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Error("message", slg.Err(err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"message": "Unauthorized"})
		}
	}
}

func RequireModerator(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.auth.RequireModerator"
		reqID := middleware.GetReqID(r.Context())
		role, ok := r.Context().Value("role").(string)

		if !ok {
			message := "failed to get role"
			log.Error(message, fn)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log = slg.WithLogger(fn, reqID)

		if role != "moderator" {
			message := "Forbidden"
			log.Error(message)
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		next.ServeHTTP(w, r)
	}
}

//func Validate(log *slog.Logger, next http.Handler) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const fn = "middleware.auth.Validate"
//		reqID := middleware.GetReqID(r.Context())
//
//		log = slg.SetupLogger(fn, reqID)
//
//		var user entity.User
//
//		err := render.DecodeJSON(r.Body, &user)
//		if err != nil {
//			message := "failed to decode"
//
//			log.Error(message)
//			render.Status(r, http.StatusBadRequest)
//			render.JSON(w, r, map[string]string{"message": message})
//		}
//
//		if !auth.IsValidEmail(user.Email) || user.Password == "" {
//			message := "invalid data"
//
//			log.Error(message)
//			render.Status(r, http.StatusBadRequest)
//			render.JSON(w, r, map[string]string{"message": message})
//		}
//
//	}
//}
