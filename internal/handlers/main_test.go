package handlers

import (
	"testing"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T, store db.Store) Handler {
	h, err := NewDefaultHandler(store)
	require.NoError(t, err)

	return h
}
