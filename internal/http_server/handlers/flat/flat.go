package flat

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/http_server/sender"
	"avito_tech/internal/lib/logger/slg"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FlatStorage
type FlatStorage interface {
	CreateF(flat entity.Flat) (int64, error)
	Update(flat entity.Flat, idMod uuid.UUID) error
	GetSubscribers(houseID int64) ([]string, error)
}

func Create(log *slog.Logger, storage FlatStorage, sender *sender.Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.flat.Create"
		reqID := middleware.GetReqID(r.Context())
		username := r.Context().Value("username").(uuid.UUID)

		log = slg.WithLogger(fn, reqID)

		var flat entity.Flat

		err := render.DecodeJSON(r.Body, &flat)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log.Info("request body decoded", slog.Any("request", reqID))

		flat.UserID = username

		id, err := storage.CreateF(flat)
		if err != nil {
			message := "failed to add flat"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log.Info("flat added", slog.Any("request", reqID))

		go func() {
			subscribers, err := storage.GetSubscribers(flat.HouseID)
			if err != nil {
				log.Error("failed to get subscribers", slg.Err(err))
				return
			}

			message := fmt.Sprintf("New flat added in house %d: Number %d, Price %d, Rooms %d, Status %s", flat.HouseID, flat.Number, flat.Price, flat.Rooms, flat.Status)

			log.Info(message)

			for _, email := range subscribers {
				err := sender.SendEmail(r.Context(), email, message)
				if err != nil {
					log.Error("failed to send email", slg.Err(err))
				}

				log.Info("message sent by email " + email)
			}
		}()

		flat.ID = id
		flat.Status = "created"
		render.JSON(w, r, flat)
	}
}

func Update(log *slog.Logger, storage FlatStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.flat.Update"
		reqID := middleware.GetReqID(r.Context())
		username := r.Context().Value("username").(uuid.UUID)

		log = slg.WithLogger(fn, reqID)

		var flat entity.Flat
		err := render.DecodeJSON(r.Body, &flat)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		err = storage.Update(flat, username)
		if err != nil {
			message := "failed to update flat"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message})
			return
		}

		log.Info("flat update", slog.Any("request", reqID))

		render.JSON(w, r, flat)
	}
}
