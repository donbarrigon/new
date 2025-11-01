package validator

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/validation"
)

// ================================================================
//                       CREATE ROLES
// ================================================================

type RoleStore struct {
	Name string `json:"name"`
}

func (r *RoleStore) Rules() validation.Rules {
	return validation.Rules{
		"name": {
			"required": {},
			"between":  {"3", "255"},
			"regex":    {"kebab-case"},
			"unique":   {"roles", "name"},
		},
	}
}

func (r *RoleStore) PrepareForValidation(c *handler.Context) *err.ValidationError {
	return err.NewValidationError()
}

func NewRoleStore(c *handler.Context) (*RoleStore, error) {
	v := &RoleStore{}
	e := validation.Body(c, v)
	return v, e
}

// ================================================================
//                       CREATE PERMISSIONS
// ================================================================

type PermissionStore struct {
	Name string `json:"name"`
}

func (p *PermissionStore) Rules() validation.Rules {
	return validation.Rules{
		"name": {
			"required": {},
			"between":  {"3", "255"},
			"regex":    {"^[a-z0-9]+(?:[ -][a-z0-9]+)*$"}, // sin mayúsculas, solo minúsculas, números, guiones y espacios
			"unique":   {"permissions", "name"},
		},
	}
}

func (p *PermissionStore) PrepareForValidation(c *handler.Context) *err.ValidationError {
	return err.NewValidationError()
}

func NewPermissionStore(c *handler.Context) (*PermissionStore, error) {
	v := &PermissionStore{}
	e := validation.Body(c, v)
	return v, e
}
