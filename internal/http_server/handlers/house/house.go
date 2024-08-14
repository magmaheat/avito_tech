package house

import (
	"avito_tech/internal/entity"
	"avito_tech/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type ResponseCreate struct {
	Message   string       `json:"message"`
	RequestID string       `json:"request_id"`
	Body      entity.House `json:"body"`
}

type ResponseFlats struct {
	Status string        `json:"status"`
	Body   []entity.Flat `json:"body"`
}

type Storage interface {
	Create(house entity.House) error
	GetFlats(idHouse int64, role string) ([]entity.Flat, error)
}

func Create(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.house.Create"
		reqID := middleware.GetReqID(r.Context())

		log = setupLogger(fn, reqID)

		var req entity.House

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		err = storage.Create(req)
		if err != nil {
			message := "failed to add house"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		message := "house added"
		log.Info(message)

		render.JSON(w, r, ResponseCreate{
			Message:   message,
			RequestID: reqID,
			Body:      req,
		})
	}
}

func Flats(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.house.Flats"
		reqID := middleware.GetReqID(r.Context())
		role := r.Context().Value("role").(string)

		log = setupLogger(fn, reqID)

		id := chi.URLParam(r, "id")
		if id == "" {
			message := "id is empty"
			log.Info(message)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		newID, _ := strconv.Atoi(id)

		var resFlats []entity.Flat
		var err error

		resFlats, err = storage.GetFlats(int64(newID), role)

		if err != nil {
			message := "failed to get flats"
			log.Error(fn, map[string]error{message: err})
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		log.Info("got flats")

		render.JSON(w, r, ResponseFlats{
			Status: "Ok",
			Body:   resFlats,
		})

	}
}

func setupLogger(fn, reqID string) *slog.Logger {
	return slog.With(
		slog.String("fn", fn),
		slog.String("id_request", reqID),
	)
}
