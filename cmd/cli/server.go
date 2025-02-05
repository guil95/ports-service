package cli

import (
	"errors"
	"github.com/guil95/ports-service/internal/core/application"
	"github.com/guil95/ports-service/internal/infra/adapters/repository"
	"github.com/guil95/ports-service/internal/infra/server/http/handler"
	"github.com/guil95/ports-service/pkg/database"
	"github.com/guil95/ports-service/pkg/graceful"
	"github.com/spf13/cobra"
	"log/slog"
	"net/http"
)

var ServeCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start HTTP server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("Starting HTTP server on :8080")

		db := database.NewPostgresDB()
		repo := repository.NewPostgresRepository(db)
		service := application.NewService(repo, nil)
		httpHandler := handler.NewHTTPHandler(service)

		server := &http.Server{
			Addr:    ":8080",
			Handler: httpHandler,
		}

		serverErr := make(chan error, 1)
		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				serverErr <- err
			}
			close(serverErr)
		}()

		ctx := graceful.WaitForShutdown()
		select {
		case <-ctx.Done():
			slog.Info("Shutting down server...")
		case err := <-serverErr:
			slog.Error("Server error", "error", err)
			return err
		}

		if err := graceful.Shutdown(server, graceful.DefaultShutdownOptions()); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
			return err
		}

		slog.Info("Server stopped gracefully")
		return nil
	},
}
