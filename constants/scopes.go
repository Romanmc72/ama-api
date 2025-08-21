package constants

import "maps"

var adminScopes = map[string]string{
	"role": "question/admin",
}

// returns a copy of the admin scopes map to prevent external modification
func GetAdminScopes() map[string]string {
	adminScopesCopy := make(map[string]string)
	maps.Copy(adminScopesCopy, adminScopes)
	return adminScopesCopy
}
