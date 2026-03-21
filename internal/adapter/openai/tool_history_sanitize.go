package openai

import (
	"regexp"
	"strings"
)

var leakedToolHistoryPattern = regexp.MustCompile(`(?is)\[TOOL_CALL_HISTORY\][\s\S]*?\[/TOOL_CALL_HISTORY\]|\[TOOL_RESULT_HISTORY\][\s\S]*?\[/TOOL_RESULT_HISTORY\]`)

func sanitizeLeakedToolHistory(text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}
	cleaned := leakedToolHistoryPattern.ReplaceAllString(text, "")
	return strings.TrimSpace(cleaned)
}
