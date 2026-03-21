package openai

import "testing"

func TestSanitizeLeakedToolHistoryRemovesMarkerBlocks(t *testing.T) {
	raw := "前缀\n[TOOL_CALL_HISTORY]\nfunction.name: exec\nfunction.arguments: {}\n[/TOOL_CALL_HISTORY]\n后缀"
	got := sanitizeLeakedToolHistory(raw)
	if got != "前缀\n\n后缀" {
		t.Fatalf("unexpected sanitized content: %q", got)
	}
}

func TestFlushToolSieveDropsToolHistoryLeak(t *testing.T) {
	var state toolStreamSieveState
	chunk := "[TOOL_CALL_HISTORY]\nstatus: already_called\nfunction.name: exec\nfunction.arguments: {}\n[/TOOL_CALL_HISTORY]"
	evts := processToolSieveChunk(&state, chunk, []string{"exec"})
	if len(evts) != 0 {
		t.Fatalf("expected no immediate output before history block is complete, got %+v", evts)
	}
	flushed := flushToolSieve(&state, []string{"exec"})
	if len(flushed) != 0 {
		t.Fatalf("expected history block to be swallowed, got %+v", flushed)
	}
}
