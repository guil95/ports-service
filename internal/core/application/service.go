package application

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/guil95/ports-service/internal/core/domain"
)

type service struct {
	repo   domain.RepositoryPort
	parser domain.ParserPort
}

const batchSize = 200

func NewService(repo domain.RepositoryPort, parser domain.ParserPort) domain.ServicePort {
	return &service{repo: repo, parser: parser}
}

func (s *service) CreateOrUpdate(ctx context.Context, port domain.Port) error {
	if len(port.Unlocs) == 0 {
		return domain.ErrInvalidPort
	}
	portID := strings.ToUpper(port.Unlocs[0])
	port.ID = &portID

	return s.repo.SaveBulk(ctx, []domain.Port{port})
}

func (s *service) ImportPorts(ctx context.Context) error {
	portCh, errCh := s.parser.Parse(ctx) // pipeline pattern
	var batch []domain.Port
	for {
		select {
		case <-ctx.Done():
			// Save items that remains on the batch
			if len(batch) > 0 {
				if err := s.repo.SaveBulk(ctx, batch); err != nil {
					return err
				}
			}

			return ctx.Err()
		case port, ok := <-portCh:
			if !ok {
				if len(batch) > 0 {
					// Save items that remains on the batch
					if err := s.repo.SaveBulk(ctx, batch); err != nil {
						return err
					}
				}
				return nil
			}

			batch = append(batch, port)

			if len(batch) == batchSize {
				// Save items as batch
				if err := s.repo.SaveBulk(ctx, batch); err != nil {
					return err
				}
				batch = nil
			}

		case err := <-errCh:
			if len(batch) > 0 {
				// Save items that remains on the batch
				if saveErr := s.repo.SaveBulk(ctx, batch); saveErr != nil {
					return fmt.Errorf("error to save batch: %v; error: %w", saveErr, err)
				}
			}
			return err
		}
	}
}

func (s *service) FindByID(ctx context.Context, portID string) (*domain.Port, error) {
	port, err := s.repo.FindByID(ctx, portID)
	if err != nil {
		slog.Error("error to get a port", "error", err)
		return nil, err
	}

	return port, nil
}
