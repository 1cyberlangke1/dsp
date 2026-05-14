package settings

import (
	"strings"

	"ds2api/internal/config"
)

func validateMergedRuntimeSettings(current config.RuntimeConfig, incoming *config.RuntimeConfig) error {
	merged := current
	if incoming != nil {
		if incoming.AccountMaxInflight > 0 {
			merged.AccountMaxInflight = incoming.AccountMaxInflight
		}
		if incoming.AccountMaxQueue > 0 {
			merged.AccountMaxQueue = incoming.AccountMaxQueue
		}
		if incoming.GlobalMaxInflight > 0 {
			merged.GlobalMaxInflight = incoming.GlobalMaxInflight
		}
		if incoming.TokenRefreshIntervalHours > 0 {
			merged.TokenRefreshIntervalHours = incoming.TokenRefreshIntervalHours
		}
		if strings.TrimSpace(incoming.AccountScheduleMode) != "" {
			merged.AccountScheduleMode = incoming.AccountScheduleMode
		}
		if incoming.AccountStickyReuseCount > 0 {
			merged.AccountStickyReuseCount = incoming.AccountStickyReuseCount
		}
	}
	return validateRuntimeSettings(merged)
}

func (h *Handler) applyRuntimeSettings() {
	if h == nil || h.Store == nil || h.Pool == nil {
		return
	}
	accountCount := len(h.Store.Accounts())
	maxPer := h.Store.RuntimeAccountMaxInflight()
	recommended := defaultRuntimeRecommended(accountCount, maxPer)
	maxQueue := h.Store.RuntimeAccountMaxQueue(recommended)
	global := h.Store.RuntimeGlobalMaxInflight(recommended)
	h.Pool.ApplyRuntimeLimits(maxPer, maxQueue, global, h.Store.RuntimeAccountScheduleMode(), h.Store.RuntimeAccountStickyReuseCount())
}

func defaultRuntimeRecommended(accountCount, maxPer int) int {
	if maxPer <= 0 {
		maxPer = 1
	}
	if accountCount <= 0 {
		return maxPer
	}
	return accountCount * maxPer
}
