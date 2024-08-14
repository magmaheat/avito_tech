package house

//
//import (
//	"github.com/go-chi/chi/v5"
//	"github.com/go-chi/chi/v5/middleware"
//	"github.com/go-chi/render"
//	"log/slog"
//	"net/http"
//)
//
//type Storage interface {
//	GetFlats(idHouse int64) pos
//}
//
//func Flats(log *slog.Logger, storage Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const fn = "handlers.house.Get"
//		reqId := middleware.GetReqID(r.Context())
//
//		log.With(
//			slog.String("fn", fn),
//			slog.String("request_id", reqId),
//		)
//
//		id := chi.URLParam(r, "id")
//		if id == "" {
//			message := "id is empty"
//			log.Info(message)
//			render.JSON(w, r, answer(message, reqId))
//			return
//		}
//
//	}
//}
