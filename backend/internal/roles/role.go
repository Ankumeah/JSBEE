package roles

import (
	"fmt"
)

// This stuct reprents a new role
// with its own permissions
//
// This stuct is not to be used
// as, rather use the vars in `roles.go`
type Role struct {
	Role      string `json:"role" binding:"required"`
	CanReview bool   `json:"can_review" binding:"required"`
	CanEditDB bool   `json:"can_edit_db" binding:"required"`
	CanBeGod  bool   `json:"can_be_god" binding:"required"`
}

// So that `Role` can be scanned by the db later
func (r *Role) Scan(value any) error {
	switch v := value.(type) {
	case string:
		switch v {
		case "owner":
			*r = Owner
		case "admin":
			*r = Admin
		case "reviewer":
			*r = Reviewer
		case "viewer":
			*r = Viewer
		default:
			return fmt.Errorf("%v is not a role", value)
		}
	default:
		return fmt.Errorf("Cannot convert %v into a role", v)
	}

	return nil
}

func newRole(
	role string,
	canReview bool,
	canEditDB bool,
	canBeGod bool,
) Role {
	return Role{
		Role:      role,
		CanReview: canReview,
		CanEditDB: canEditDB,
		CanBeGod:  canBeGod,
	}
}
