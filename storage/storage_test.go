package storage_test

import (
	"testing"

	"github.com/gourses/demo/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestInvalidConnString(t *testing.T) {
	s, err := storage.New("not valid conn str")
	assert.Nil(t, s)
	assert.Error(t, err)
}
