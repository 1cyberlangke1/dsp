package client

import (
	"strings"
	"testing"
)

func TestDeviceIDDeterministic(t *testing.T) {
	id1 := DeviceID("test@example.com")
	id2 := DeviceID("test@example.com")
	if id1 != id2 {
		t.Fatalf("expected deterministic device id, got %q vs %q", id1, id2)
	}
}

func TestDeviceIDDifferentWhenAccountDiffers(t *testing.T) {
	id1 := DeviceID("a@example.com")
	id2 := DeviceID("b@example.com")
	if id1 == id2 {
		t.Fatal("expected different device ids for different accounts")
	}
}

func TestDeviceIDFormat(t *testing.T) {
	id := DeviceID("test@example.com")
	// SHA512 → 64 bytes → Base64, so 64*4/3 = 85.3 → ceil 88 with padding
	if len(id) != 88 {
		t.Fatalf("expected device id length 88 (64 bytes base64), got %d: %q", len(id), id)
	}
	// Must be valid Base64 charset
	for _, r := range id {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", r) {
			t.Fatalf("invalid base64 character %q in device id", r)
		}
	}
}

func TestDeviceIDNotEmpty(t *testing.T) {
	if DeviceID("") == "" {
		t.Fatal("expected non-empty device id even for empty input")
	}
}
