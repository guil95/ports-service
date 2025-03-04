package cli

import (
	"context"
	"fmt"
	"github.com/guil95/ports-service/database"
	"github.com/guil95/ports-service/graceful"
	"log/slog"
	"os"

	"github.com/guil95/ports-service/internal/core/application"
	"github.com/guil95/ports-service/internal/infra/adapters/parser"
	"github.com/guil95/ports-service/internal/infra/adapters/repository"
	"github.com/spf13/cobra"
)

func init() {
	ImportCmd.Flags().StringP("file", "f", "", "Path to JSON file (required)")
	_ = ImportCmd.MarkFlagRequired("file")
}

var ImportCmd = &cobra.Command{
	Use:          "import",
	Short:        "Import ports from a JSON file",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := graceful.WaitForShutdown()
		filePath, _ := cmd.Flags().GetString("file")

		slog.Info("Starting import process", "file", filePath)

		result := make(chan error, 1)
		go func() {
			defer close(result)
			result <- runImport(ctx, filePath)
		}()

		select {
		case err := <-result:
			if err != nil {
				slog.Error("Import failed", "error", err)
				return err
			}
			slog.Info("Import completed successfully")
			return nil
		case <-ctx.Done():
			slog.Info("Received shutdown signal, stopping import...")
		}

		err := <-result
		if err != nil {
			slog.Error("Import interrupted with error", "error", err)
			return fmt.Errorf("import interrupted with error: %w", err)
		}
		slog.Info("Import stopped gracefully")
		return nil
	},
}

func runImport(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	_, err = file.Stat()
	if err != nil {
		return fmt.Errorf("could not obtain stat, handle error: %w", err)
	}

	db := database.NewPostgresDB()
	repo := repository.NewPostgresRepository(db)
	jsonParser := parser.NewJSONParser(file)
	service := application.NewService(repo, jsonParser)

	if err := service.ImportPorts(ctx); err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	fmt.Println("Ports imported successfully!")
	return nil
}
