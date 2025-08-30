package constants

import "maps"

var adminScopes = map[string]any{
	"role": "question/admin",
}

// returns a copy of the admin scopes map to prevent external modification
func GetAdminScopes() map[string]any {
	adminScopesCopy := make(map[string]any)
	maps.Copy(adminScopesCopy, adminScopes)
	return adminScopesCopy
}
