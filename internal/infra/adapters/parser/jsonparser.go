package parser

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/guil95/ports-service/internal/core/domain"
)

type jsonParser struct {
	reader io.Reader
}

func NewJSONParser(reader io.Reader) domain.ParserPort {
	return &jsonParser{
		reader,
	}
}

func (p *jsonParser) Parse(ctx context.Context) (<-chan domain.Port, <-chan error) {
	portCh := make(chan domain.Port)
	errCh := make(chan error, 1)

	go func() {
		defer close(portCh)
		defer close(errCh)

		decoder := json.NewDecoder(p.reader)

		token, err := decoder.Token()
		if err != nil {
			slog.ErrorContext(ctx, "error to read initial token", "error", err)
			errCh <- domain.ErrInvalidJson
			return
		}
		if delim, ok := token.(json.Delim); !ok || delim != '{' {
			slog.ErrorContext(ctx, "error to read initial token")
			errCh <- domain.ErrInvalidJson
			return
		}

		for decoder.More() {
			// reading keys to read line by line
			keyToken, err := decoder.Token()
			if err != nil {
				slog.ErrorContext(ctx, "error to read key")
				errCh <- domain.ErrInvalidJson
				return
			}

			key, ok := keyToken.(string)
			if !ok {
				slog.ErrorContext(ctx, "key should be string")
				errCh <- domain.ErrInvalidJson
				return
			}

			var port domain.Port
			if err := decoder.Decode(&port); err != nil {
				slog.ErrorContext(ctx, "error to decode value to the key", "key", key, "error", err)
				errCh <- domain.ErrInvalidJson
				return
			}

			port.ID = &key
			portCh <- port
		}

		token, err = decoder.Token()
		if err != nil {
			slog.ErrorContext(ctx, "error to read last token", "error", err)
			errCh <- domain.ErrInvalidJson
			return
		}
		if delim, ok := token.(json.Delim); !ok || delim != '}' {
			errCh <- domain.ErrInvalidJson
			return
		}
	}()

	return portCh, errCh
}
