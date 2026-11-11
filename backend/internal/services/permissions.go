package services

import "slices"

// RequiredPermissions are the GW2 API key permissions the app needs to work.
// Kept as a fixed list so new ArenaNet permissions don't invalidate existing keys
var RequiredPermissions = []string{"account", "characters", "inventories", "wallet"}

// MissingPermissions returns the required permissions not in granted, sorted
func MissingPermissions(granted []string) []string {
	var missing []string
	for _, permission := range RequiredPermissions {
		if !slices.Contains(granted, permission) {
			missing = append(missing, permission)
		}
	}
	slices.Sort(missing)
	return missing
}
