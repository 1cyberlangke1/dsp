package protocol

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestSharedConstantsLoaded(t *testing.T) {
	cfg := sharedConstants{}
	if err := json.Unmarshal(sharedConstantsJSON, &cfg); err != nil {
		t.Fatalf("failed to parse shared constants: %v", err)
	}
	client := normalizeClientConstants(cfg.Client)
	if ClientVersion != client.Version {
		t.Fatalf("unexpected client version=%q", ClientVersion)
	}
	wantUserAgent := client.Name + "/" + client.Version + " Android/" + client.AndroidAPILevel
	if BaseHeaders["User-Agent"] != wantUserAgent {
		t.Fatalf("unexpected user agent=%q", BaseHeaders["User-Agent"])
	}
	if BaseHeaders["x-client-platform"] != "android" {
		t.Fatalf("unexpected base header x-client-platform=%q", BaseHeaders["x-client-platform"])
	}
	if BaseHeaders["x-client-version"] != ClientVersion {
		t.Fatalf("unexpected base header x-client-version=%q", BaseHeaders["x-client-version"])
	}
	if BaseHeaders["Content-Type"] != "application/json" {
		t.Fatalf("unexpected base header Content-Type=%q", BaseHeaders["Content-Type"])
	}
	if len(SkipContainsPatterns) == 0 {
		t.Fatal("expected skip contains patterns to be loaded")
	}
	if _, ok := SkipExactPathSet["response/search_status"]; !ok {
		t.Fatal("expected response/search_status in exact skip path set")
	}
}

func TestClientHeadersDerivedFromSharedVersion(t *testing.T) {
	client := normalizeClientConstants(clientConstants{
		Name:            "DeepSeek",
		Platform:        "android",
		Version:         "9.8.7",
		AndroidAPILevel: "28",
		Locale:          "zh_CN",
	})
	headers := buildBaseHeaders(client, map[string]string{
		"User-Agent":       "stale",
		"x-client-version": "stale",
	})
	if headers["User-Agent"] != "DeepSeek/9.8.7 Android/28" {
		t.Fatalf("unexpected derived user agent=%q", headers["User-Agent"])
	}
	if headers["x-client-version"] != "9.8.7" {
		t.Fatalf("unexpected derived client version=%q", headers["x-client-version"])
	}
}

func TestBaseHeadersIncludeDeepSeekBrandHeader(t *testing.T) {
	if BaseHeaders["x-client-bundle-id"] != "com.deepseek.chat" {
		t.Fatalf("expected x-client-bundle-id=com.deepseek.chat, got %q", BaseHeaders["x-client-bundle-id"])
	}
}

func TestBaseHeadersIncludeAcceptEncoding(t *testing.T) {
	if BaseHeaders["Accept-Encoding"] != "gzip" {
		t.Fatalf("expected Accept-Encoding=gzip, got %q", BaseHeaders["Accept-Encoding"])
	}
}

func TestRangersIDIs19Digits(t *testing.T) {
	if RangersID == "" {
		t.Fatal("expected RangersID to be non-empty")
	}
	if len(RangersID) != 19 {
		t.Fatalf("expected RangersID to be exactly 19 digits, got %d: %q", len(RangersID), RangersID)
	}
	for _, r := range RangersID {
		if r < '0' || r > '9' {
			t.Fatalf("RangersID must be all digits, got %q", RangersID)
		}
	}
}

func TestRangersIDInBaseHeaders(t *testing.T) {
	if BaseHeaders["x-rangers-id"] != RangersID {
		t.Fatalf("expected BaseHeaders[x-rangers-id]=%q, got %q", RangersID, BaseHeaders["x-rangers-id"])
	}
}

func TestTimezoneOffsetInBaseHeaders(t *testing.T) {
	v, ok := BaseHeaders["x-client-timezone-offset"]
	if !ok || v == "" {
		t.Fatalf("expected x-client-timezone-offset in base headers")
	}
	offset, err := strconv.Atoi(v)
	if err != nil {
		t.Fatalf("x-client-timezone-offset must be integer, got %q: %v", v, err)
	}
	if offset != 28800 {
		t.Fatalf("expected fixed timezone offset 28800, got %d", offset)
	}
}

func TestBuildBaseHeadersIncludesDynamicFields(t *testing.T) {
	client := normalizeClientConstants(clientConstants{
		Name:    "DeepSeek",
		Version: "2.0.0",
	})
	headers := buildBaseHeaders(client, nil)
	if headers["x-client-platform"] != "android" {
		t.Fatalf("expected x-client-platform=android, got %q", headers["x-client-platform"])
	}
	if headers["x-client-locale"] != "zh_CN" {
		t.Fatalf("expected x-client-locale=zh_CN, got %q", headers["x-client-locale"])
	}
	if _, ok := headers["x-client-timezone-offset"]; !ok {
		t.Fatal("expected x-client-timezone-offset")
	}
	if _, ok := headers["x-rangers-id"]; !ok {
		t.Fatal("expected x-rangers-id")
	}
}

func TestStaticBaseHeaders(t *testing.T) {
	if BaseHeaders["Host"] != "chat.deepseek.com" {
		t.Fatalf("expected Host header, got %q", BaseHeaders["Host"])
	}
	if BaseHeaders["Accept"] != "application/json" {
		t.Fatalf("expected Accept header, got %q", BaseHeaders["Accept"])
	}
	if BaseHeaders["accept-charset"] != "UTF-8" {
		t.Fatalf("expected accept-charset header, got %q", BaseHeaders["accept-charset"])
	}
}

func TestRangersIDRemainsStableDuringProcess(t *testing.T) {
	if RangersID == "" {
		t.Fatal("RangersID should be set during init")
	}
	id1 := RangersID
	id2 := RangersID
	if id1 != id2 {
		t.Fatal("RangersID changed between reads")
	}
}

func TestRangersIDForSeedIsStableAndDistinct(t *testing.T) {
	id1 := RangersIDForSeed("acct-a")
	id2 := RangersIDForSeed("acct-a")
	id3 := RangersIDForSeed("acct-b")
	if id1 != id2 {
		t.Fatalf("expected stable rangers id, got %q vs %q", id1, id2)
	}
	if id1 == id3 {
		t.Fatalf("expected distinct rangers ids for different seeds, both were %q", id1)
	}
	if len(id1) != 19 || len(id3) != 19 {
		t.Fatalf("expected 19-digit ids, got %q and %q", id1, id3)
	}
}

func TestUserAgentContainsAppNameVersionAndOS(t *testing.T) {
	ua := BaseHeaders["User-Agent"]
	if !strings.HasPrefix(ua, "DeepSeek/") {
		t.Fatalf("User-Agent should start with DeepSeek/, got %q", ua)
	}
	if !strings.Contains(ua, "Android/") {
		t.Fatalf("User-Agent should contain Android/, got %q", ua)
	}
	parts := strings.Split(ua, " ")
	// Expected format: "DeepSeek/X.Y.Z Android/N"
	if len(parts) < 2 {
		t.Fatalf("User-Agent should have at least 2 parts, got %q", ua)
	}
	if !strings.HasPrefix(parts[1], "Android/") {
		t.Fatalf("User-Agent second part should be Android/, got %q", parts[1])
	}
	// Verify the Android API level parses as int
	androidPart := strings.TrimPrefix(parts[1], "Android/")
	if _, err := strconv.Atoi(androidPart); err != nil {
		t.Fatalf("Android API level should be numeric, got %q: %v", androidPart, err)
	}
}

func TestUserAgentUsesAndroid28(t *testing.T) {
	if got := BaseHeaders["User-Agent"]; !strings.Contains(got, "Android/28") {
		t.Fatalf("expected Android/28 in User-Agent, got %q", got)
	}
}
