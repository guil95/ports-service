package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownOptions struct {
	Timeout time.Duration
}

func DefaultShutdownOptions() ShutdownOptions {
	return ShutdownOptions{
		Timeout: 10 * time.Second,
	}
}

func WaitForShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		fmt.Println("\nReceived shutdown signal")
		cancel()
	}()

	return ctx
}

func Shutdown(server *http.Server, opts ShutdownOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	fmt.Println("Server stopped gracefully")
	return nil
}
