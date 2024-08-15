package flat

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Storage interface {
	CreateFlat(flat entity.Flat) error
}

func Create(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.flat.Create"
		reqID := middleware.GetReqID(r.Context())

		log = slg.SetupLogger(fn, reqID)

		var flat entity.Flat

		err := render.DecodeJSON(r.Body, &flat)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"Error": message})
			return
		}

		log.Info("request body decoded", slog.Any("request", reqID))

		err = storage.CreateFlat(flat)
		if err != nil {
			message := "failed to add flat"
			log.Error(message)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"Error": message})
			return
		}

		log.Info("flat added", slog.Any("request", reqID))

		flat.Status = "created"
		render.JSON(w, r, flat)
	}
}
