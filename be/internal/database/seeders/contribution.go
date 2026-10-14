package seeders

import "fmt"

// Contribution is what a module adds to the seed run: permission catalog rows and
// extra seeders (run after the core RBAC and user seeders).
type Contribution struct {
	Permissions []PermissionDef
	Seeders     []Seeder
}

// MergePermissions returns the core catalog followed by every contribution's
// permissions, rejecting duplicate keys or IDs.
func MergePermissions(contribs ...Contribution) ([]PermissionDef, error) {
	all := append([]PermissionDef{}, DefaultPermissions()...)
	for _, c := range contribs {
		all = append(all, c.Permissions...)
	}

	keys := make(map[string]bool, len(all))
	ids := make(map[string]bool, len(all))
	for _, p := range all {
		if keys[p.Key] {
			return nil, fmt.Errorf("seeders: duplicate permission key %q", p.Key)
		}
		if ids[p.ID] {
			return nil, fmt.Errorf("seeders: duplicate permission id %q (%s)", p.ID, p.Key)
		}
		keys[p.Key] = true
		ids[p.ID] = true
	}
	return all, nil
}
