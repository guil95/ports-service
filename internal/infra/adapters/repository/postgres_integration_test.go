//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/guil95/ports-service/tests/suite"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository(t *testing.T) {
	t.Run("save bulk of ports and find by id", func(t *testing.T) {
		ctx := context.Background()
		postgresContainer, db := suite.SetupPostgresContainer(t)
		defer postgresContainer.Terminate(ctx)
		defer db.Close()

		repo := NewPostgresRepository(db)

		ports := []domain.Port{
			{
				ID:          stringPtr("CNCGU"),
				Name:        "China",
				City:        "Changshu",
				Country:     "China",
				Alias:       []string{"Zhangjiagang", "Suzhou", "Taicang"},
				Regions:     []string{"Region1", "Region2"},
				Coordinates: []float64{120.752503, 31.653686},
				Province:    "Jiangsu",
				Timezone:    "Asia/Shanghai",
				Unlocs:      []string{"CNCGU"},
				Code:        "57076",
			},
			{
				ID:          stringPtr("CNBJO"),
				Name:        "Beijiao",
				City:        "Beijiao",
				Country:     "China",
				Alias:       []string{},
				Regions:     []string{"Region1", "Region2"},
				Coordinates: []float64{119.92, 26.35},
				Province:    "Fujian",
				Timezone:    "Asia/Shanghai",
				Unlocs:      []string{"CNBJO"},
				Code:        "57016",
			},
		}

		err := repo.SaveBulk(context.Background(), ports)
		assert.NoError(t, err)

		retrievedPort, err := repo.FindByID(ctx, *ports[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, ports[0], *retrievedPort)

		retrievedPort2, err := repo.FindByID(ctx, *ports[1].ID)
		assert.NoError(t, err)
		assert.Equal(t, ports[1], *retrievedPort2)

		nonExistentPort, err := repo.FindByID(ctx, "NON_EXISTENT")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPortNotFound)
		assert.Nil(t, nonExistentPort)
	})
}

func stringPtr(s string) *string {
	return &s
}
