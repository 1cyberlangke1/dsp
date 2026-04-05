package util

import (
	"strings"
	"testing"
)

func TestBuildToolCallInstructions_ExecCommandUsesCmdExample(t *testing.T) {
	out := BuildToolCallInstructions([]string{"exec_command"})
	if !strings.Contains(out, `<tool_name>exec_command</tool_name>`) {
		t.Fatalf("expected exec_command in examples, got: %s", out)
	}
	if !strings.Contains(out, `<parameters>{"cmd":"pwd"}</parameters>`) {
		t.Fatalf("expected cmd parameter example for exec_command, got: %s", out)
	}
}

func TestBuildToolCallInstructions_ExecuteCommandUsesCommandExample(t *testing.T) {
	out := BuildToolCallInstructions([]string{"execute_command"})
	if !strings.Contains(out, `<tool_name>execute_command</tool_name>`) {
		t.Fatalf("expected execute_command in examples, got: %s", out)
	}
	if !strings.Contains(out, `<parameters>{"command":"pwd"}</parameters>`) {
		t.Fatalf("expected command parameter example for execute_command, got: %s", out)
	}
}

func TestBuildToolCallInstructions_IncludesFailureRecoveryAndNoFabricationRules(t *testing.T) {
	out := BuildToolCallInstructions([]string{"WebSearch"})
	if !strings.Contains(out, "Never claim that a tool was run") {
		t.Fatalf("expected no-fabrication guard, got: %s", out)
	}
	if !strings.Contains(out, "If a tool call fails, is rejected, or returns no usable result") {
		t.Fatalf("expected failure recovery rule, got: %s", out)
	}
	if !strings.Contains(out, "Never output internal tool-planning traces") {
		t.Fatalf("expected no internal trace rule, got: %s", out)
	}
}

func TestBuildToolCallInstructions_AllowsTextBeforeXMLButNotAfter(t *testing.T) {
	out := BuildToolCallInstructions([]string{"read_file"})
	if !strings.Contains(out, "Any explanatory text must appear before that block, never after it.") {
		t.Fatalf("expected clarified xml placement rule, got: %s", out)
	}
	if strings.Contains(out, "No text before, no text after") {
		t.Fatalf("expected contradictory no-text-before wording removed, got: %s", out)
	}
}
