package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wlcmtunknwndth/test_ozon/graph"
	"github.com/wlcmtunknwndth/test_ozon/internal/auth"
	"github.com/wlcmtunknwndth/test_ozon/internal/config"
	"github.com/wlcmtunknwndth/test_ozon/internal/storage/inmemory"
	"github.com/wlcmtunknwndth/test_ozon/internal/storage/postgres"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"net/http"
)

func main() {

	cfg := config.MustLoad()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	var authService auth.Auth
	var graphQL graph.Resolver
	if cfg.UseDB {
		pg, err := postgres.New(&cfg.DB)
		if err != nil {
			slog.Error("couldn't run storage", slogAttr.SlogErr("main", err))
			return
		}
		slog.Info("initialized postgres")
		authService = auth.Auth{Db: pg}
		graphQL = graph.Resolver{Storage: pg}
	} else {
		inMemoryStorage := inmemory.New()
		authService = auth.Auth{Db: inMemoryStorage}
		graphQL = graph.Resolver{Storage: inMemoryStorage}

	}
	router.Use(authService.MiddlewareAuth())
	router.Post("/login", authService.LogIn)
	router.Post("/register", authService.Register)
	router.Post("/logout", authService.LogOut)
	router.Post("/delete_user", authService.DeleteUser)

	gql := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graphQL}))

	router.Handle("/query", gql)

	slog.Info("connect to server for GraphQL playground", slogAttr.SlogInfo("address", cfg.Server.Address))

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to run server", slogAttr.SlogErr("main", err))
	}
	slog.Info("server closed")
}
