package main

import (
	"avito_tech/internal/config"
	"avito_tech/internal/http_server/handlers/auth"
	"avito_tech/internal/http_server/handlers/flat"
	"avito_tech/internal/http_server/handlers/house"
	mdr "avito_tech/internal/http_server/middleware/auth"
	"avito_tech/internal/lib/logger/slg"
	"avito_tech/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	os.Setenv("MY_SIGNING_KEY", "pussy")
	os.Setenv("CONFIG_PATH", "../../config/local.yaml")
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/dummyLogin", auth.DummyLogin(log, storage))
	router.Post("/flat/create", mdr.JWTAuth(log, flat.Create(log, storage)))
	router.Post("/flat/update", mdr.JWTAuth(log, flat.Update(log, storage)))
	//router.Post("/login")
	//router.Post("/register")
	router.Post("/house/create", mdr.JWTAuth(log, mdr.RequireModerator(log, house.Create(log, storage))))
	router.Get("/house/{id}", mdr.JWTAuth(log, house.Flats(log, storage)))

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

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// /dummyLogin get +
// /house/create post +
// /house/{id} get +

// /flat/create post
// /flat/update post

// /register post
// /login post

// /house/{id}/subscribe post
