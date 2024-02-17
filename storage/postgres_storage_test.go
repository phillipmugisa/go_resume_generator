package storage

import (
	"testing"
)

func TestDB(t *testing.T) {
	_, err := NewPostgresStorage()
	if err != nil {
		t.Error(err)
	}
}
