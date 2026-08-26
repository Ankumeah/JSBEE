package roles

// This stuct reprents a new role
// with its own permissions
//
// This stuct is not to be used
// as, rather use the vars in `roles.go`
type Role struct {
	Value     string `json:"value" binding:"required"`
	CanReview bool   `json:"can_review" binding:"required"`
	CanEditDB bool   `json:"can_edit_db" binding:"required"`
	CanBeGod  bool   `json:"can_be_god" binding:"required"`
}

func newRole(
	value string,
	canReview bool,
	canEditDB bool,
	canBeGod bool,
) Role {
	return Role{
		Value:     value,
		CanReview: canReview,
		CanEditDB: canEditDB,
		CanBeGod:  canBeGod,
	}
}
