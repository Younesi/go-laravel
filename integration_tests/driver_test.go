//go:build integration

package integration_tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/younesi/atlas"
)

func TestOpenDB_Success(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port("5432/tcp")),
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Get the mapped port for PostgreSQL
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := "postgres://testuser:testpassword@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	atlas := &atlas.Atlas{
		InfoLog: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	db, err := atlas.OpenDB("postgres", dsn)
	require.NoError(t, err)

	// Verify the connection is successful
	require.NoError(t, db.Ping())

	// Close the connection
	require.NoError(t, db.Close())
}
