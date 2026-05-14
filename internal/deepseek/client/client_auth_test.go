package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"ds2api/internal/auth"
	"ds2api/internal/config"
	"ds2api/internal/deepseek/protocol"
)

type bodyCaptureRoundTripper struct {
	capturedBody []byte
}

func (t *bodyCaptureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	t.capturedBody = buf.Bytes()
	req.Body.Close()
	// Return a valid login response
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}

func TestLoginPayloadStructureForEmail(t *testing.T) {
	rt := &bodyCaptureRoundTripper{}
	client := &Client{
		regular:   doerFunc(func(r *http.Request) (*http.Response, error) { return rt.RoundTrip(r) }),
		fallback:  &http.Client{Transport: rt},
		maxRetries: 3,
	}
	_, _ = client.Login(context.Background(), config.Account{
		Email:    "user@example.com",
		Password: "secret123",
	})

	var payload map[string]any
	if err := json.Unmarshal(rt.capturedBody, &payload); err != nil {
		t.Fatalf("failed to decode login payload: %v", err)
	}

	if payload["email"] != "user@example.com" {
		t.Fatalf("expected email=user@example.com, got %#v", payload["email"])
	}
	if payload["password"] != "secret123" {
		t.Fatalf("expected password=secret123, got %#v", payload["password"])
	}
	if payload["os"] != "Android" {
		t.Fatalf("expected os=Android, got %#v", payload["os"])
	}
	deviceID, _ := payload["device_id"].(string)
	if deviceID == "" || len(deviceID) != 88 {
		t.Fatalf("expected device_id to be 88-char base64, got %q (len=%d)", deviceID, len(deviceID))
	}
	// Email login: mobile and area_code should be empty string
	if payload["mobile"] != "" {
		t.Fatalf("expected mobile=\"\" for email login, got %#v", payload["mobile"])
	}
	if payload["area_code"] != "" {
		t.Fatalf("expected area_code=\"\" for email login, got %#v", payload["area_code"])
	}
}

func TestLoginPayloadStructureForMobile(t *testing.T) {
	rt := &bodyCaptureRoundTripper{}
	client := &Client{
		regular:   doerFunc(func(r *http.Request) (*http.Response, error) { return rt.RoundTrip(r) }),
		fallback:  &http.Client{Transport: rt},
		maxRetries: 3,
	}
	_, _ = client.Login(context.Background(), config.Account{
		Mobile:   "13800138000",
		Password: "secret123",
	})

	var payload map[string]any
	if err := json.Unmarshal(rt.capturedBody, &payload); err != nil {
		t.Fatalf("failed to decode login payload: %v", err)
	}

	if payload["password"] != "secret123" {
		t.Fatalf("expected password=secret123, got %#v", payload["password"])
	}
	if payload["os"] != "Android" {
		t.Fatalf("expected os=Android, got %#v", payload["os"])
	}
	// Mobile login: email should be nil
	if payload["email"] != nil {
		t.Fatalf("expected email=nil for mobile login, got %#v", payload["email"])
	}
	mobile, _ := payload["mobile"].(string)
	if mobile == "" {
		t.Fatalf("expected non-empty mobile for mobile login")
	}
	// area_code is nil when mobile has no country prefix
	if payload["area_code"] != nil {
		t.Fatalf("expected area_code=nil for plain mobile, got %#v", payload["area_code"])
	}
}

func TestLoginSendsHeaders(t *testing.T) {
	var capturedHeaders http.Header
	client := &Client{
		regular: doerFunc(func(r *http.Request) (*http.Response, error) {
			capturedHeaders = r.Header
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}, nil
		}),
		fallback:  &http.Client{},
		maxRetries: 3,
	}
	_, _ = client.Login(context.Background(), config.Account{
		Email:    "user@example.com",
		Password: "secret123",
	})

	if capturedHeaders.Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type=application/json, got %q", capturedHeaders.Get("Content-Type"))
	}
	// Login should not send Authorization (token not yet acquired)
	if capturedHeaders.Get("authorization") != "" {
		t.Fatalf("expected no authorization header on login, got %q", capturedHeaders.Get("authorization"))
	}
}

func TestCreateSessionSendsPayload(t *testing.T) {
	var capturedBody []byte
	client := &Client{
		regular: doerFunc(func(r *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			capturedBody = buf.Bytes()
			r.Body.Close()
			// Return a valid session response
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}, nil
		}),
		fallback:  &http.Client{},
		maxRetries: 3,
	}
	_, _ = client.CreateSession(context.Background(), &auth.RequestAuth{DeepSeekToken: "test-token"}, 3)

	var payload map[string]any
	if err := json.Unmarshal(capturedBody, &payload); err != nil {
		t.Fatalf("failed to decode create session payload: %v", err)
	}
	if payload["agent"] != "chat" {
		t.Fatalf("expected agent=chat, got %#v", payload["agent"])
	}
}

func TestExtractCreateSessionIDSupportsLegacyShape(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"biz_data": map[string]any{
				"id": "legacy-session-id",
			},
		},
	}

	if got := extractCreateSessionID(resp); got != "legacy-session-id" {
		t.Fatalf("expected legacy session id, got %q", got)
	}
}

func TestExtractCreateSessionIDSupportsNestedChatSessionShape(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"biz_data": map[string]any{
				"chat_session": map[string]any{
					"id":         "nested-session-id",
					"model_type": "default",
				},
			},
		},
	}

	if got := extractCreateSessionID(resp); got != "nested-session-id" {
		t.Fatalf("expected nested session id, got %q", got)
	}
}

func TestAuthHeadersIncludeBearerAndBaseHeaders(t *testing.T) {
	// Save and restore BaseHeaders to not affect other tests
	saved := protocol.BaseHeaders
	defer func() { protocol.BaseHeaders = saved }()

	protocol.BaseHeaders = map[string]string{
		"x-client-platform": "android",
		"x-client-version":  "2.1.0",
	}

	client := &Client{}
	headers := client.authHeaders("my-token")

	if headers["authorization"] != "Bearer my-token" {
		t.Fatalf("expected authorization=Bearer my-token, got %q", headers["authorization"])
	}
	if headers["x-client-platform"] != "android" {
		t.Fatalf("expected x-client-platform=android, got %q", headers["x-client-platform"])
	}
	if headers["x-client-version"] != "2.1.0" {
		t.Fatalf("expected x-client-version=2.1.0, got %q", headers["x-client-version"])
	}
}

func TestLoginDeviceIDIsDeterministicPerAccount(t *testing.T) {
	emailPayloads := make([]map[string]any, 2)
	for i := 0; i < 2; i++ {
		rt := &bodyCaptureRoundTripper{}
		client := &Client{
			regular:   doerFunc(func(r *http.Request) (*http.Response, error) { return rt.RoundTrip(r) }),
			fallback:  &http.Client{Transport: rt},
			maxRetries: 3,
		}
		_, _ = client.Login(context.Background(), config.Account{
			Email:    "same@example.com",
			Password: "secret123",
		})
		var p map[string]any
		json.Unmarshal(rt.capturedBody, &p)
		emailPayloads[i] = p
	}
	id1, _ := emailPayloads[0]["device_id"].(string)
	id2, _ := emailPayloads[1]["device_id"].(string)
	if id1 != id2 {
		t.Fatalf("expected same device_id for same account, got %q vs %q", id1, id2)
	}
}
