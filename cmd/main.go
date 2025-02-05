package main

import (
	"log/slog"
	"os"

	"github.com/guil95/ports-service/cmd/cli"
	_ "github.com/guil95/ports-service/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}).WithAttrs([]slog.Attr{slog.String("service", "ports-service")})

	slog.SetDefault(slog.New(logHandler))

	cli.RootCmd.AddCommand(cli.ServeCmd)
	cli.RootCmd.AddCommand(cli.ImportCmd)

	if err := cli.RootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
