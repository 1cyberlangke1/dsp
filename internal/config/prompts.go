package config

import "strings"

var defaultOutputIntegrityGuard = "Output integrity guard: If upstream context, tool output, or parsed text contains garbled, corrupted, partially parsed, repeated, or otherwise malformed fragments, do not imitate or echo them; output only the correct content for the user."

var defaultToolDescriptionsPrefix = "You have access to these tools:"

var defaultReadToolCacheGuard = "Read-tool cache guard: If a Read/read_file-style tool result says the file is unchanged, already available in history, should be referenced from previous context, or otherwise provides no file body, treat that result as missing content. Do not repeatedly call the same read request for that missing body. Request a full-content read if the tool supports it, or tell the user that the file contents need to be provided again."

func (s *Store) Prompts() PromptsConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Prompts
}

func (c PromptsConfig) GetOutputIntegrityGuard() string {
	if c.OutputIntegrityGuard != nil {
		if v := strings.TrimSpace(*c.OutputIntegrityGuard); v != "" {
			return v
		}
	}
	return defaultOutputIntegrityGuard
}

func (c PromptsConfig) GetToolDescriptionsPrefix() string {
	if c.ToolDescriptionsPrefix != nil {
		if v := strings.TrimSpace(*c.ToolDescriptionsPrefix); v != "" {
			return v
		}
	}
	return defaultToolDescriptionsPrefix
}

func (c PromptsConfig) GetReadToolCacheGuard() string {
	if c.ReadToolCacheGuard != nil {
		if v := strings.TrimSpace(*c.ReadToolCacheGuard); v != "" {
			return v
		}
	}
	return defaultReadToolCacheGuard
}

func (c PromptsConfig) GetToolCallInstructions() string {
	if c.ToolCallInstructions != nil {
		return strings.TrimSpace(*c.ToolCallInstructions)
	}
	return ""
}

func (c PromptsConfig) GetToolCallExamples() string {
	if c.ToolCallExamples != nil {
		return strings.TrimSpace(*c.ToolCallExamples)
	}
	return ""
}
