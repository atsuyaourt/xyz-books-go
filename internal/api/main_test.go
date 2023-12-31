package api

import (
	"testing"

	db "github.com/emiliogozo/xyz-books/db/sqlc"
	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		GinMode:     gin.TestMode,
		APIBasePath: "/api/v1",
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
