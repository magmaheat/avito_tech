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

type ResponseCreateHouse struct {
	Message   string       `json:"message"`
	RequestID string       `json:"request_id"`
	House     entity.House `json:"house"`
}

type ResponseGetFlats struct {
	Status string        `json:"status"`
	Flat   []entity.Flat `json:"flat"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=HouseStorage
type HouseStorage interface {
	CreateH(house entity.House) (int64, error)
	GetFlats(idHouse int64, role string) ([]entity.Flat, error)
}

func Create(log *slog.Logger, storage HouseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.house.Create"
		reqID := middleware.GetReqID(r.Context())

		log = slg.SetupLogger(fn, reqID)

		var req entity.House

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		id, err := storage.CreateH(req)
		if err != nil {
			message := "failed to add house"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"message": message, "request_id": reqID})
			return
		}

		req.ID = id
		message := "house added"
		log.Info(message)

		render.JSON(w, r, ResponseCreateHouse{
			Message:   message,
			RequestID: reqID,
			House:     req,
		})
	}
}

func Flats(log *slog.Logger, storage HouseStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.house.Flats"
		reqID := middleware.GetReqID(r.Context())
		role := r.Context().Value("role").(string)

		log = slg.SetupLogger(fn, reqID)

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

		render.JSON(w, r, ResponseGetFlats{
			Status: "Ok",
			Flat:   resFlats,
		})
	}
}
