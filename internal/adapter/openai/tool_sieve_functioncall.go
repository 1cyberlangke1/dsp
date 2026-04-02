package openai

import "strings"

func findQuotedFunctionCallKeyStart(s string) int {
	lower := strings.ToLower(s)
	quotedIdx := findFunctionCallKeyStart(lower, `"functioncall"`)
	baretIdx := findFunctionCallKeyStart(lower, "functioncall")

	switch {
	case quotedIdx < 0:
		return baretIdx
	case baretIdx < 0:
		return quotedIdx
	case quotedIdx < baretIdx:
		return quotedIdx
	default:
		return baretIdx
	}
}

func findFunctionCallKeyStart(lower, key string) int {
	for from := 0; from < len(lower); {
		rel := strings.Index(lower[from:], key)
		if rel < 0 {
			return -1
		}
		idx := from + rel
		if !hasJSONObjectContextPrefix(lower[:idx]) {
			from = idx + 1
			continue
		}
		if !hasJSONKeyBoundary(lower, idx, len(key)) {
			from = idx + 1
			continue
		}
		j := idx + len(key)
		for j < len(lower) && (lower[j] == ' ' || lower[j] == '\t' || lower[j] == '\r' || lower[j] == '\n') {
			j++
		}
		if j < len(lower) && lower[j] == ':' {
			return idx
		}
		from = idx + 1
	}
	return -1
}

func hasJSONObjectContextPrefix(prefix string) bool {
	return strings.LastIndex(prefix, "{") >= 0
}

func hasJSONKeyBoundary(s string, idx, keyLen int) bool {
	if idx > 0 {
		prev := s[idx-1]
		if isLowerAlphaNumeric(prev) {
			return false
		}
	}
	if end := idx + keyLen; end < len(s) {
		next := s[end]
		if isLowerAlphaNumeric(next) {
			return false
		}
	}
	return true
}

func isLowerAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '_'
}
