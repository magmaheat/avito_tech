package main

import (
	"avito_tech/internal/config"
	"avito_tech/internal/http_server/handlers/auth"
	"avito_tech/internal/http_server/handlers/flat"
	"avito_tech/internal/http_server/handlers/house"
	mdr "avito_tech/internal/http_server/middleware/auth"
	send "avito_tech/internal/http_server/sender"
	"avito_tech/internal/lib/logger/slg"
	"avito_tech/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	//os.Setenv("MY_SIGNING_KEY", "pussy")
	//os.Setenv("CONFIG_PATH", "../../config/local.yaml")
	cfg := config.MustLoad()

	log := slg.SetupLogger(cfg.Env)

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1)
	}

	sender := send.New()
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/dummyLogin", auth.DummyLogin(log, storage))
	router.Post("/login", auth.Login(log, storage))
	router.Post("/register", auth.Register(log, storage))

	router.Post("/house/create", mdr.JWTAuth(log, mdr.RequireModerator(log, house.Create(log, storage))))
	router.Get("/house/{id}", mdr.JWTAuth(log, house.GetAllFlats(log, storage)))
	router.Post("/house/{id}/subscribe", house.Subscribe(log, storage))

	router.Post("/flat/create", mdr.JWTAuth(log, flat.Create(log, storage, sender)))
	router.Post("/flat/update", mdr.JWTAuth(log, mdr.RequireModerator(log, flat.Update(log, storage))))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
