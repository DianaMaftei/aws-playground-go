package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../test.env")
	if err != nil {
		panic(fmt.Errorf("unable to load .env file: %w", err))
	}

	ctx := context.Background()

	seedDataPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}
	mountPath := seedDataPath + "/../scripts/"

	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": os.Getenv("DB_PASSWORD"),
			"MYSQL_DATABASE":      os.Getenv("DB_NAME"),
		},
		Mounts: testcontainers.Mounts(testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: mountPath,
			},
			Target: "/docker-entrypoint-initdb.d",
		}),
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create a generic container: %w", err))
	}

	defer container.Terminate(ctx)

	p, _ := container.MappedPort(ctx, "3306")
	os.Setenv("DB_PORT", p.Port())

	Init()

	exitVal := m.Run()
	os.Exit(exitVal)
}
