package webui

import "embed"

// embeddedAdminFS bundles the built admin webui so serverless runtime does not
// depend on function includeFiles behavior.
//
//go:embed assets/admin
var embeddedAdminFS embed.FS
