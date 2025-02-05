package domain

import (
	"context"
)

// ServicePort (Primary Port)
type ServicePort interface {
	CreateOrUpdate(ctx context.Context, port Port) error
	FindByID(ctx context.Context, portID string) (*Port, error)
	ImportPorts(ctx context.Context) error
}

// RepositoryPort (Secondary Port)
type RepositoryPort interface {
	SaveBulk(ctx context.Context, port []Port) error
	FindByID(ctx context.Context, id string) (*Port, error)
}

// ParserPort (Secondary Port)
type ParserPort interface {
	Parse(ctx context.Context) (<-chan Port, <-chan error)
}
