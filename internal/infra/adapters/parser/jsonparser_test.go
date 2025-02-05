//go:build unit

package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/guil95/ports-service/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONParser(t *testing.T) {
	t.Run("test with valid json should return success", func(t *testing.T) {
		jsonData := `{
		"AEAJM": {
			"name": "Ajman",
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
		parser := NewJSONParser(reader)

		portCh, errCh := parser.Parse(context.Background())

		var ports []domain.Port
		for port := range portCh {
			ports = append(ports, port)
		}

		require.Len(t, ports, 1)
		assert.Equal(t, "AEAJM", *ports[0].ID)
		assert.Equal(t, "Ajman", ports[0].Name)
		assert.Equal(t, "Ajman", ports[0].City)
		assert.Equal(t, "United Arab Emirates", ports[0].Country)
		assert.Equal(t, []float64{55.5136433, 25.4052165}, ports[0].Coordinates)
		assert.Equal(t, "Ajman", ports[0].Province)
		assert.Equal(t, "Asia/Dubai", ports[0].Timezone)
		assert.Equal(t, []string{"AEAJM"}, ports[0].Unlocs)
		assert.Equal(t, "52000", ports[0].Code)

		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		default:
		}
	})
	t.Run("test with invalid json (json without delim in the end), should return error", func(t *testing.T) {
		jsonData := `{
		"AEAJM": {
			"name": "Ajman",
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
	`
		reader := strings.NewReader(jsonData)
		parser := NewJSONParser(reader)

		portCh, errCh := parser.Parse(context.Background())

		select {
		case err := <-errCh:
			assert.True(t, errors.Is(err, domain.ErrInvalidJson))
		case <-portCh:
			t.Fatal("expected an error but got a port")
		}
	})
	t.Run("test with invalid type data should return error", func(t *testing.T) {
		jsonData := `{
		"AEAJM": {
			"name": "Ajman",
			"city": "Ajman",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": "invalid", // Invalid coordinates
			"province": "Ajman",
			"timezone": "Asia/Dubai",
			"unlocs": [
				"AEAJM"
			],
			"code": "52000"
		}
	}`

		reader := strings.NewReader(jsonData)
		parser := NewJSONParser(reader)

		portCh, errCh := parser.Parse(context.Background())

		select {
		case err := <-errCh:
			assert.True(t, errors.Is(err, domain.ErrInvalidJson))
		case <-portCh:
			t.Fatal("expected an error but got a port")
		}
	})
}
