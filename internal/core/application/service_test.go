//go:build unit

package application

import (
	"context"
	"errors"
	"testing"

	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/guil95/ports-service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	t.Run("create port with success", func(t *testing.T) {
		repoMock := new(mocks.RepositoryPort)
		parserMock := new(mocks.ParserPort)
		portsService := NewService(repoMock, parserMock)
		ctx := context.Background()

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

		portToSave := port
		portToSave.ID = &port.Unlocs[0]
		repoMock.On("SaveBulk", ctx, []domain.Port{portToSave}).Return(nil)

		err := portsService.CreateOrUpdate(ctx, port)
		assert.NoError(t, err)
	})

	t.Run("create port with error", func(t *testing.T) {
		repoMock := new(mocks.RepositoryPort)
		parserMock := new(mocks.ParserPort)
		portsService := NewService(repoMock, parserMock)
		ctx := context.Background()

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

		repoError := errors.New("internal error")
		portToSave := port
		portToSave.ID = &port.Unlocs[0]
		repoMock.On("SaveBulk", ctx, []domain.Port{portToSave}).Return(repoError)

		err := portsService.CreateOrUpdate(ctx, port)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repoError)
	})

	t.Run("find port by ID", func(t *testing.T) {
		repoMock := new(mocks.RepositoryPort)
		parserMock := new(mocks.ParserPort)
		portsService := NewService(repoMock, parserMock)
		ctx := context.Background()
		id := "CNCGU"
		portResponse := &domain.Port{
			ID:          &id,
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

		repoMock.On("FindByID", ctx, id).Return(portResponse, nil)

		p, err := portsService.FindByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, portResponse, p)
	})

	t.Run("find port by ID return not found", func(t *testing.T) {
		repoMock := new(mocks.RepositoryPort)
		parserMock := new(mocks.ParserPort)
		portsService := NewService(repoMock, parserMock)
		ctx := context.Background()
		id := "CNCGU"

		repoMock.On("FindByID", ctx, id).Return(nil, domain.ErrPortNotFound)

		p, err := portsService.FindByID(ctx, id)
		assert.Nil(t, p)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPortNotFound)
	})
}

func TestImports(t *testing.T) {
	t.Run("import with success", func(t *testing.T) {
		importData := []domain.Port{{
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
		}, {
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
		}
		ctx := context.Background()
		repoMock := &mocks.RepositoryPort{}
		parserMock := &mocks.ParserPort{}

		portCh := make(chan domain.Port, len(importData))
		errCh := make(chan error, 1)

		for _, port := range importData {
			portCh <- port
		}
		close(portCh)
		errCh <- nil
		close(errCh)

		batchToSave := importData
		batchToSave[0].ID = importData[0].ID
		batchToSave[1].ID = importData[1].ID

		parserMock.On("Parse", ctx).Return((<-chan domain.Port)(portCh), (<-chan error)(errCh))
		repoMock.On("SaveBulk", ctx, batchToSave).Return(nil).Times(1)

		service := NewService(repoMock, parserMock)
		err := service.ImportPorts(ctx)
		assert.NoError(t, err)
	})
	t.Run("import with error", func(t *testing.T) {
		importData := []domain.Port{{
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
		}, {
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
		}
		ctx := context.Background()
		repoMock := &mocks.RepositoryPort{}
		parserMock := &mocks.ParserPort{}

		portCh := make(chan domain.Port, len(importData))
		errCh := make(chan error, 1)
		parseError := errors.New("internal error")
		for _, port := range importData {
			portCh <- port
		}
		close(portCh)
		errCh <- parseError
		close(errCh)

		batchToSave := importData
		batchToSave[0].ID = importData[0].ID
		batchToSave[1].ID = importData[1].ID

		parserMock.On("Parse", ctx).Return((<-chan domain.Port)(portCh), (<-chan error)(errCh))
		repoMock.On("SaveBulk", ctx, batchToSave).Return(nil).Times(1)

		s := NewService(repoMock, parserMock)
		err := s.ImportPorts(ctx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, parseError)
	})
}
