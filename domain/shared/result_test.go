package shared_test

import (
	"testing"

	"github.com/i-nishimura/goatodo/domain/shared"
)

func TestOk(t *testing.T) {
	t.Run("Ok result holds value and reports success", func(t *testing.T) {
		result := shared.Ok(42)

		if !result.IsOk() {
			t.Error("expected IsOk() to be true")
		}
		if result.IsErr() {
			t.Error("expected IsErr() to be false")
		}
		if result.Value() != 42 {
			t.Errorf("expected Value() to be 42, got %d", result.Value())
		}
	})
}

func TestErr(t *testing.T) {
	t.Run("Err result holds error and reports failure", func(t *testing.T) {
		result := shared.Err[int]("something went wrong")

		if result.IsOk() {
			t.Error("expected IsOk() to be false")
		}
		if !result.IsErr() {
			t.Error("expected IsErr() to be true")
		}
		if result.Error() != "something went wrong" {
			t.Errorf("expected Error() to be 'something went wrong', got '%s'", result.Error())
		}
	})
}
