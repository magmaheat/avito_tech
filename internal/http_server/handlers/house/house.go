package house

import (
	"avito_tech/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Address   string `json:"address"`
	Year      int64  `json:"year"`
	Developer string `json:"developer"`
	UserType  string `json:"user_type"`
}

type Response struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type House interface {
	Create(address string, year int64, developer string) error
}

func Create(log *slog.Logger, house House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.house.Create"
		reqId := middleware.GetReqID(r.Context())

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", reqId),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			message := "failed to decode request body"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Message:   message,
				RequestID: reqId,
			})
			return
		}

		err = house.Create(req.Address, req.Year, req.Developer)
		if err != nil {
			message := "failed to add house"
			log.Error(message, slg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Message:   message,
				RequestID: reqId,
			})
			return
		}

		message := "house added"
		log.Info(message)

		render.JSON(w, r, Response{
			Message:   message,
			RequestID: reqId,
		})

	}
}
