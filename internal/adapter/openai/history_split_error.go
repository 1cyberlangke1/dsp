package openai

import (
	"net/http"

	"ds2api/internal/deepseek"
)

func mapHistorySplitError(err error) (int, string) {
	switch {
	case deepseek.IsManagedUnauthorizedError(err):
		return http.StatusUnauthorized, "Account token is invalid. Please re-login the account in admin."
	case deepseek.IsDirectUnauthorizedError(err):
		return http.StatusUnauthorized, "Invalid token. If this should be a DS2API key, add it to config.keys first."
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
