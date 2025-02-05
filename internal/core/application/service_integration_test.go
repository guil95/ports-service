//go:build integration

package application

import (
	"context"
	"github.com/guil95/ports-service/internal/infra/adapters/parser"
	"strings"
	"testing"

	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/guil95/ports-service/internal/infra/adapters/repository"
	"github.com/guil95/ports-service/tests/suite"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationService(t *testing.T) {
	t.Run("test create and find with success", func(t *testing.T) {
		ctx := context.Background()
		container, db := suite.SetupPostgresContainer(t)
		defer container.Terminate(ctx)

		repo := repository.NewPostgresRepository(db)
		s := NewService(repo, nil)
		port := domain.Port{
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
		}
		err := s.CreateOrUpdate(ctx, port)
		assert.NoError(t, err)
		portResponse, err := s.FindByID(ctx, port.Unlocs[0])
		assert.NoError(t, err)
		assert.Equal(t, port.Name, portResponse.Name)
	})
	t.Run("test imports with success", func(t *testing.T) {
		ctx := context.Background()
		container, db := suite.SetupPostgresContainer(t)
		defer container.Terminate(ctx)

		repo := repository.NewPostgresRepository(db)
		jsonData := `{
			"AEAJM": {
				"name": "TEST 1",
				"city": "Ajman",
				"country": "United Arab Emirates",
				"alias": [],
				"regions": [],
				"coordinates": [
					55.5136433,
					25.4052165
				],
				"province": "Ajman",
				"timezone": "Asia/Dubai",
				"unlocs": [
					"AEAJM"
				],
				"code": "52000"
			},
			"ABCDE": {
				"name": "TEST 2",
				"city": "Ajman",
				"country": "United Arab Emirates",
				"alias": [],
				"regions": [],
				"coordinates": [
					55.5136433,
					25.4052165
				],
				"province": "Ajman",
				"timezone": "Asia/Dubai",
				"unlocs": [
					"AEAJM"
				],
				"code": "52000"
			}
		}`

		reader := strings.NewReader(jsonData)
		p := parser.NewJSONParser(reader)
		s := NewService(repo, p)

		err := s.ImportPorts(ctx)
		assert.NoError(t, err)

		portResponse, err := s.FindByID(ctx, "AEAJM")
		assert.NoError(t, err)
		assert.Equal(t, "TEST 1", portResponse.Name)

		portResponse2, err := s.FindByID(ctx, "ABCDE")
		assert.NoError(t, err)
		assert.Equal(t, "TEST 2", portResponse2.Name)

		_, err = s.FindByID(ctx, "asassa")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPortNotFound)
	})
}
