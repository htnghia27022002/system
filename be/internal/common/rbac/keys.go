package rbac

import "strings"

const (
	ActionView   = "view"
	ActionModify = "modify"
)

func Key(resource, action string) string {
	return resource + ":" + action
}

// Allowed reports whether granted permissions satisfy required.
// Supports legacy CRUD keys (read/create/update/delete) for transitional data.
func Allowed(granted []string, required string) bool {
	for _, item := range granted {
		if item == required {
			return true
		}
	}

	parts := strings.SplitN(required, ":", 2)
	if len(parts) != 2 {
		return false
	}

	resource, action := parts[0], parts[1]
	switch action {
	case ActionView:
		for _, legacy := range []string{Key(resource, "read")} {
			for _, item := range granted {
				if item == legacy {
					return true
				}
			}
		}
	case ActionModify:
		for _, legacy := range []string{
			Key(resource, "create"),
			Key(resource, "update"),
			Key(resource, "delete"),
		} {
			for _, item := range granted {
				if item == legacy {
					return true
				}
			}
		}
	}

	return false
}
