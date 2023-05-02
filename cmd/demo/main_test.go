//go:build integration

package main

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_main(t *testing.T) {
	// Setup App
	setupTestEnv(t)
	go main()
	time.Sleep(time.Second)

	// Query notes API
	resp, err := http.Get("http://127.0.0.1:8080/notes")
	require.NoError(t, err)
	defer func() { assert.NoError(t, resp.Body.Close()) }()

	// Validate response
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "[]\n", string(data))
}

func setupTestEnv(t *testing.T) {
	compose, err := tc.NewDockerCompose("../../docker-compose.yaml")
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal))
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	require.NoError(t, compose.WaitForService("db", wait.ForHealthCheck()).Up(ctx, tc.Wait(true)))
}
