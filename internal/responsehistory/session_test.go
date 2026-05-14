package responsehistory

import (
	"net/http/httptest"
	"testing"
)

func TestShouldCaptureIgnoresLegacyStreamQueryFlags(t *testing.T) {
	t.Run("prepare flag", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/chat/completions?__stream_prepare=1", nil)
		if !shouldCapture(req) {
			t.Fatal("expected legacy prepare flag to be ignored")
		}
	})

	t.Run("release flag", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/chat/completions?__stream_release=1", nil)
		if !shouldCapture(req) {
			t.Fatal("expected legacy release flag to be ignored")
		}
	})
}
