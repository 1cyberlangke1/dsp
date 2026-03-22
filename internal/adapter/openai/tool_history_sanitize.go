package openai

import (
	"regexp"
)

var leakedToolHistoryPattern = regexp.MustCompile(`(?is)\[TOOL_CALL_HISTORY\][\s\S]*?\[/TOOL_CALL_HISTORY\]|\[TOOL_RESULT_HISTORY\][\s\S]*?\[/TOOL_RESULT_HISTORY\]`)
var emptyJSONFencePattern = regexp.MustCompile("(?is)```json\\s*```")
var leakedToolCallArrayPattern = regexp.MustCompile(`(?is)\[\{\s*"function"\s*:\s*\{[\s\S]*?\}\s*,\s*"id"\s*:\s*"call[^"]*"\s*,\s*"type"\s*:\s*"function"\s*}\]`)
var leakedToolResultBlobPattern = regexp.MustCompile(`(?is)<\s*\|\s*tool\s*\|\s*>\s*\{[\s\S]*?"tool_call_id"\s*:\s*"call[^"]*"\s*}`)
var leakedMetaMarkerPattern = regexp.MustCompile(`(?is)<\s*\|\s*(?:assistant|tool|end_of_sentence|end_of_thinking)\s*\|\s*>`)

func sanitizeLeakedToolHistory(text string) string {
	if text == "" {
		return text
	}
	out := leakedToolHistoryPattern.ReplaceAllString(text, "")
	out = emptyJSONFencePattern.ReplaceAllString(out, "")
	out = leakedToolCallArrayPattern.ReplaceAllString(out, "")
	out = leakedToolResultBlobPattern.ReplaceAllString(out, "")
	out = leakedMetaMarkerPattern.ReplaceAllString(out, "")
	return out
}
