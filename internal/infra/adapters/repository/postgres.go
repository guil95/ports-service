package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) domain.RepositoryPort {
	return &postgresRepository{db}
}

func (r *postgresRepository) SaveBulk(ctx context.Context, ports []domain.Port) error {
	type PortDB struct {
		ID          *string     `db:"id"`
		Name        string      `db:"name"`
		City        string      `db:"city"`
		Country     string      `db:"country"`
		Alias       interface{} `db:"alias"`
		Regions     interface{} `db:"regions"`
		Coordinates interface{} `db:"coordinates"`
		Province    string      `db:"province"`
		Timezone    string      `db:"timezone"`
		Unlocs      interface{} `db:"unlocs"`
		Code        string      `db:"code"`
	}

	var portsDB []PortDB
	for _, p := range ports {
		portsDB = append(portsDB, PortDB{
			ID:          p.ID,
			Name:        p.Name,
			City:        p.City,
			Country:     p.Country,
			Alias:       pq.Array(p.Alias),
			Regions:     pq.Array(p.Regions),
			Coordinates: pq.Array(p.Coordinates),
			Province:    p.Province,
			Timezone:    p.Timezone,
			Unlocs:      pq.Array(p.Unlocs),
			Code:        p.Code,
		})
	}

	query := `
		INSERT INTO ports (id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code)
		VALUES (:id, :name, :city, :country, :alias, :regions, :coordinates, :province, :timezone, :unlocs, :code)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			city = EXCLUDED.city,
			country = EXCLUDED.country,
			alias = EXCLUDED.alias,
			regions = EXCLUDED.regions,
			coordinates = EXCLUDED.coordinates,
			province = EXCLUDED.province,
			timezone = EXCLUDED.timezone,
			unlocs = EXCLUDED.unlocs,
			code = EXCLUDED.code;
	`

	_, err := r.db.NamedExecContext(ctx, query, portsDB)
	return err
}

func (r *postgresRepository) FindByID(ctx context.Context, id string) (*domain.Port, error) {
	query := `
	SELECT id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code
	FROM ports
	WHERE LOWER(id) = LOWER($1)
	`

	var rawPort struct {
		ID          string          `db:"id"`
		Name        string          `db:"name"`
		City        string          `db:"city"`
		Country     string          `db:"country"`
		Alias       pq.StringArray  `db:"alias"`
		Regions     pq.StringArray  `db:"regions"`
		Coordinates pq.Float64Array `db:"coordinates"`
		Province    string          `db:"province"`
		Timezone    string          `db:"timezone"`
		Unlocs      pq.StringArray  `db:"unlocs"`
		Code        string          `db:"code"`
	}

	err := r.db.GetContext(ctx, &rawPort, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPortNotFound
		}
		return nil, fmt.Errorf("error fetching port: %v", err)
	}

	port := &domain.Port{
		ID:          &rawPort.ID,
		Name:        rawPort.Name,
		City:        rawPort.City,
		Country:     rawPort.Country,
		Alias:       []string(rawPort.Alias),
		Regions:     []string(rawPort.Regions),
		Coordinates: []float64(rawPort.Coordinates),
		Province:    rawPort.Province,
		Timezone:    rawPort.Timezone,
		Unlocs:      []string(rawPort.Unlocs),
		Code:        rawPort.Code,
	}

	return port, nil
}
