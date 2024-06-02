package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wlcmtunknwndth/test_ozon/graph"
	"github.com/wlcmtunknwndth/test_ozon/internal/auth"
	"github.com/wlcmtunknwndth/test_ozon/internal/config"
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

	var pg *postgres.Storage
	if cfg.UseDB {
		var err error
		pg, err = postgres.New(&cfg.DB)
		if err != nil {
			slog.Error("couldn't run storage", slogAttr.SlogErr("main", err))
			return
		}
	}

	authService := auth.Auth{Db: pg}
	router.Use(authService.MiddlewareAuth())
	router.Post("/login", authService.LogIn)
	router.Post("/register", authService.Register)
	router.Post("/logout", authService.LogOut)
	router.Post("/delete_user", authService.DeleteUser)

	gql := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: pg}}))

	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	////http.Handle("/query", srv)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gql)
	srv := http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	slog.Info("connect to server for GraphQL playground", slogAttr.SlogInfo("address", cfg.Server.Address))
	//log.Fatal(http.ListenAndServe(cfg.Server.Address, nil))
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to run server", slogAttr.SlogErr("main", err))
	}
	slog.Info("server closed")
}