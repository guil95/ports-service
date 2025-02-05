//go:build integration

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guil95/ports-service/internal/core/application"
	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/guil95/ports-service/internal/infra/adapters/repository"
	"github.com/guil95/ports-service/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Run("create and get port with success", func(t *testing.T) {
		ctx := context.Background()
		container, db := suite.SetupPostgresContainer(t)
		defer container.Terminate(ctx)

		repo := repository.NewPostgresRepository(db)
		service := application.NewService(repo, nil)
		h := NewHTTPHandler(service)

		payload := `{
		"name": "Dubai",
		"coordinates": [55.27, 25.25],
		"city": "Dubai",
		"province": "Dubayy [Dubai]",
		"country": "United Arab Emirates",
		"alias": [],
		"regions": [],
		"timezone": "Asia/Dubai",
		"unlocs": ["AEDXB"],
		"code": "52005"
	}`

		req, err := http.NewRequest(http.MethodPost, "/ports", strings.NewReader(payload))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		getReq, err := http.NewRequest(http.MethodGet, "/ports/AEDXB", nil)
		require.NoError(t, err)

		getRR := httptest.NewRecorder()
		h.ServeHTTP(getRR, getReq)

		assert.Equal(t, http.StatusOK, getRR.Code)

		var responsePort domain.Port
		err = json.NewDecoder(getRR.Body).Decode(&responsePort)
		require.NoError(t, err, "invalid response body")

		assert.Equal(t, "AEDXB", responsePort.Unlocs[0])
		assert.Equal(t, "AEDXB", *responsePort.ID)
		assert.Equal(t, "Dubai", responsePort.Name)
		assert.Equal(t, "Asia/Dubai", responsePort.Timezone)
	})
}
