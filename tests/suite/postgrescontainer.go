package suite

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgresContainer(t *testing.T) (testcontainers.Container, *sqlx.DB) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(60 * time.Second),
	}

	time.Sleep(5 * time.Second)
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		time.Sleep(5 * time.Second)
		postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	}
	require.NoError(t, err)

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	time.Sleep(5 * time.Second)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		time.Sleep(5 * time.Second)
		db, err = sqlx.Connect("postgres", connStr)
	}
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ports (
			id          VARCHAR(50) PRIMARY KEY,
			name        VARCHAR(100),
			city        VARCHAR(100),
			country     VARCHAR(100),
			alias       TEXT[],
			regions     TEXT[],
			coordinates FLOAT[],
			province    VARCHAR(100),
			timezone    VARCHAR(50),
			unlocs      TEXT[],
			code        VARCHAR(10)
		);
	`)
	require.NoError(t, err)

	return postgresContainer, db
}
