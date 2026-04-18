package openai

import (
	"strings"

	"ds2api/internal/toolcall"
)

func processToolSieveChunk(state *toolStreamSieveState, chunk string, toolNames []string) []toolStreamEvent {
	if state == nil {
		return nil
	}
	if chunk != "" {
		state.pending.WriteString(chunk)
	}
	events := make([]toolStreamEvent, 0, 2)
	if len(state.pendingToolCalls) > 0 {
		events = append(events, toolStreamEvent{ToolCalls: state.pendingToolCalls})
		state.pendingToolRaw = ""
		state.pendingToolCalls = nil
	}

	for {
		if state.capturing {
			if state.pending.Len() > 0 {
				state.capture.WriteString(state.pending.String())
				state.pending.Reset()
			}
			prefix, calls, suffix, ready := consumeToolCapture(state, toolNames)
			if !ready {
				break
			}
			captured := state.capture.String()
			state.capture.Reset()
			state.capturing = false
			state.resetIncrementalToolState()
			if len(calls) > 0 {
				if prefix != "" {
					state.noteText(prefix)
					events = append(events, toolStreamEvent{Content: prefix})
				}
				if suffix != "" {
					state.pending.WriteString(suffix)
				}
				_ = captured
				state.pendingToolCalls = calls
				continue
			}
			if prefix != "" {
				state.noteText(prefix)
				events = append(events, toolStreamEvent{Content: prefix})
			}
			if suffix != "" {
				state.pending.WriteString(suffix)
			}
			continue
		}

		pending := state.pending.String()
		if pending == "" {
			break
		}
		start := findToolSegmentStart(pending)
		if start >= 0 {
			prefix := pending[:start]
			if prefix != "" {
				state.noteText(prefix)
				events = append(events, toolStreamEvent{Content: prefix})
			}
			state.pending.Reset()
			state.capture.WriteString(pending[start:])
			state.capturing = true
			state.resetIncrementalToolState()
			continue
		}

		safe, hold := splitSafeContentForToolDetection(pending)
		if safe == "" {
			break
		}
		state.pending.Reset()
		state.pending.WriteString(hold)
		state.noteText(safe)
		events = append(events, toolStreamEvent{Content: safe})
	}

	return events
}

func flushToolSieve(state *toolStreamSieveState, toolNames []string) []toolStreamEvent {
	if state == nil {
		return nil
	}
	events := processToolSieveChunk(state, "", toolNames)
	if len(state.pendingToolCalls) > 0 {
		events = append(events, toolStreamEvent{ToolCalls: state.pendingToolCalls})
		state.pendingToolRaw = ""
		state.pendingToolCalls = nil
	}
	if state.capturing {
		consumedPrefix, consumedCalls, consumedSuffix, ready := consumeToolCapture(state, toolNames)
		if ready {
			if consumedPrefix != "" {
				state.noteText(consumedPrefix)
				events = append(events, toolStreamEvent{Content: consumedPrefix})
			}
			if len(consumedCalls) > 0 {
				events = append(events, toolStreamEvent{ToolCalls: consumedCalls})
			}
			if consumedSuffix != "" {
				state.noteText(consumedSuffix)
				events = append(events, toolStreamEvent{Content: consumedSuffix})
			}
		} else {
			content := state.capture.String()
			if content != "" {
				// If the captured text looks like an incomplete XML tool call block,
				// swallow it to prevent leaking raw XML tags to the client.
				if hasOpenXMLToolTag(content) {
					// Drop it silently — incomplete tool call.
				} else {
					state.noteText(content)
					events = append(events, toolStreamEvent{Content: content})
				}
			}
		}
		state.capture.Reset()
		state.capturing = false
		state.resetIncrementalToolState()
	}
	if state.pending.Len() > 0 {
		content := state.pending.String()
		// Safety: if pending contains XML tool tag fragments (e.g. "tool_calls>"
		// from a split closing tag), swallow them instead of leaking.
		if hasOpenXMLToolTag(content) || looksLikeXMLToolTagFragment(content) {
			// Drop it — likely an incomplete tool call fragment.
		} else {
			state.noteText(content)
			events = append(events, toolStreamEvent{Content: content})
		}
		state.pending.Reset()
	}
	return events
}

func splitSafeContentForToolDetection(s string) (safe, hold string) {
	if s == "" {
		return "", ""
	}
	if xmlIdx := findPartialXMLToolTagStart(s); xmlIdx >= 0 {
		if xmlIdx > 0 {
			return s[:xmlIdx], s[xmlIdx:]
		}
		return "", s
	}
	return s, ""
}

func findToolSegmentStart(s string) int {
	if s == "" {
		return -1
	}
	lower := strings.ToLower(s)
	bestKeyIdx := -1
	for _, tag := range xmlToolTagsToDetect {
		idx := strings.Index(lower, tag)
		if idx >= 0 && (bestKeyIdx < 0 || idx < bestKeyIdx) {
			bestKeyIdx = idx
		}
	}
	return bestKeyIdx
}

func consumeToolCapture(state *toolStreamSieveState, toolNames []string) (prefix string, calls []toolcall.ParsedToolCall, suffix string, ready bool) {
	captured := state.capture.String()
	if captured == "" {
		return "", nil, "", false
	}

	// XML tool call extraction only.
	if xmlPrefix, xmlCalls, xmlSuffix, xmlReady := consumeXMLToolCapture(captured, toolNames); xmlReady {
		return xmlPrefix, xmlCalls, xmlSuffix, true
	}
	// If XML tags are present but block is incomplete, keep buffering.
	if hasOpenXMLToolTag(captured) {
		return "", nil, "", false
	}
	return "", nil, "", false
}
