package util

import "strings"

func ResolveThinkingEnabled(req map[string]any, defaultEnabled bool) bool {
	if enabled, ok := parseThinkingSetting(req["thinking"]); ok {
		return enabled
	}
	if extraBody, ok := req["extra_body"].(map[string]any); ok {
		if enabled, ok := parseThinkingSetting(extraBody["thinking"]); ok {
			return enabled
		}
	}
	if enabled, ok := parseReasoningEffort(req["reasoning_effort"]); ok {
		return enabled
	}
	return defaultEnabled
}

func parseThinkingSetting(raw any) (bool, bool) {
	switch v := raw.(type) {
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "enabled":
			return true, true
		case "disabled":
			return false, true
		default:
			return false, false
		}
	case map[string]any:
		if typ, ok := v["type"]; ok {
			return parseThinkingSetting(typ)
		}
	}
	return false, false
}

func parseReasoningEffort(raw any) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(toString(raw))) {
	case "low", "medium", "high", "xhigh":
		return true, true
	default:
		return false, false
	}
}

func toString(raw any) string {
	if s, ok := raw.(string); ok {
		return s
	}
	return ""
}
