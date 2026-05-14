package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOkHttpTransportAddsPreemptiveHeader(t *testing.T) {
	var captured *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := New(0)
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	if captured == nil {
		t.Fatal("expected request to be captured")
	}
	if captured.Header.Get("OkHttp-Preemptive") != "1" {
		t.Fatalf("expected OkHttp-Preemptive: 1, got %q", captured.Header.Get("OkHttp-Preemptive"))
	}
}

func TestOkHttpTransportPreservesExistingHeader(t *testing.T) {
	var capturedHeaders http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := New(0)
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	req.Header.Set("OkHttp-Preemptive", "0")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	// okhttpTransport should not override if already set
	if capturedHeaders.Get("OkHttp-Preemptive") != "0" {
		t.Fatalf("expected OkHttp-Preemptive to remain 0, got %q", capturedHeaders.Get("OkHttp-Preemptive"))
	}
}

func TestNewFallbackClientHasOkHttpHeader(t *testing.T) {
	var captured *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewFallbackClient(0, nil)
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	if captured == nil {
		t.Fatal("expected request to be captured")
	}
	if captured.Header.Get("OkHttp-Preemptive") != "1" {
		t.Fatalf("expected OkHttp-Preemptive: 1 on fallback client, got %q", captured.Header.Get("OkHttp-Preemptive"))
	}
}

func TestClientReusesCookieJar(t *testing.T) {
	requestCount := 0
	var cookieHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		cookieHeader = r.Header.Get("Cookie")
		if requestCount == 1 {
			// Set a cookie on first request
			w.Header().Set("Set-Cookie", "ds_session_id=sess-123; Path=/")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := New(0)

	// First request — should get the cookie
	req1, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	resp1, _ := client.Do(req1)
	resp1.Body.Close()

	// Second request — should carry the cookie from jar
	req2, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	resp2, _ := client.Do(req2)
	resp2.Body.Close()

	if requestCount < 2 {
		t.Fatal("expected two requests to be made")
	}
	if cookieHeader == "" {
		t.Fatal("expected Cookie header on second request via CookieJar")
	}
}
