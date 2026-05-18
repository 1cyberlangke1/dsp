package settings

import (
	"ds2api/internal/prompt"
	"ds2api/internal/promptcompat"
	"ds2api/internal/toolcall"
)

func (h *Handler) ApplyPromptsSettings() {
	if h == nil || h.Store == nil {
		return
	}
	p := h.Store.Prompts()

	prompt.SetOutputIntegrityGuard(p.GetOutputIntegrityGuard())

	toolcall.SetCustomInstructions(p.GetToolCallInstructions())
	toolcall.SetCustomExamples(p.GetToolCallExamples())

	promptcompat.SetToolDescriptionsPrefix(p.GetToolDescriptionsPrefix())
	promptcompat.SetReadToolCacheGuard(p.GetReadToolCacheGuard())
}
