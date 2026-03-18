package shared_test

import (
	"testing"

	"github.com/i-nishimura/goatodo/domain/shared"
)

func TestNewID(t *testing.T) {
	t.Run("NewID generates a non-empty unique ID", func(t *testing.T) {
		id := shared.NewID()

		if id == "" {
			t.Error("expected non-empty ID")
		}
	})

	t.Run("NewID generates unique IDs each call", func(t *testing.T) {
		id1 := shared.NewID()
		id2 := shared.NewID()

		if id1 == id2 {
			t.Errorf("expected unique IDs, got %s and %s", id1, id2)
		}
	})
}
