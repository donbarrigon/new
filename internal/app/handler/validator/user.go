package validator

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/validation"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
//                             STORE
// ================================================================

type UserStore struct {
	Email                string        `json:"email"`
	Password             string        `json:"password"`
	PasswordConfirmation string        `json:"passwordConfirmation"`
	Nickname             string        `json:"nickname"`
	Name                 string        `json:"name"`
	Phone                string        `json:"phone,omitempty"`
	Discord              string        `json:"discord,omitempty"`
	CityID               bson.ObjectID `json:"cityId"`
}

func (u *UserStore) Rules() validation.Rules {
	return validation.Rules{
		"email": {
			"required": {},
			"regex":    {"email"},
			"between":  {"3", "254"},
			"unique":   {"users", "email"},
		},
		"password": {
			"required": {},
			"between":  {"8", "32"},
		},
		"passwordConfirmation": {
			"required": {},
		},
		"nickname": {
			"required": {},
			"between":  {"3", "255"},
			"regex":    {"nickname"},
			"unique":   {"users", "profile.nickname"},
		},
		"name": {
			"required": {},
			"between":  {"3", "255"},
		},
		"phone": {
			"regex": {"phone"},
		},
		"discord": {
			"regex":   {"discord"},
			"between": {"3", "32"},
		},
		"cityId": {
			"required": {},
			"exists":   {"cities", "_id"},
		},
	}
}
func (u *UserStore) PrepareForValidation(c *handler.Context) *err.ValidationError {
	e := err.NewValidationError()
	if u.Password != u.PasswordConfirmation {
		e.AppendM("password", "Las contraseñas no coinciden")
	}
	return e
}

func NewUserStrore(c *handler.Context) (*UserStore, error) {
	v := &UserStore{}
	e := validation.Body(c, v)
	return v, e
}

// ================================================================
//                       UPDATE PROFILE
// ================================================================

type UserUpdateProfile struct {
	Nickname string        `json:"nickname"`
	Name     string        `json:"name"`
	Phone    string        `json:"phone,omitempty"`
	Discord  string        `json:"discord,omitempty"`
	CityID   bson.ObjectID `json:"cityId"`
}

func (u *UserUpdateProfile) Rules() validation.Rules {
	return validation.Rules{
		"nickname": {
			"between": {"3", "255"},
		},
		"name": {
			"between": {"3", "255"},
		},
		"phone": {
			"regex": {"phone"},
		},
		"discord": {
			"regex":   {"discord"},
			"between": {"3", "32"},
		},
		"cityId": {
			"required": {},
			"exists":   {"cities", "_id"},
		},
	}
}

func (u *UserUpdateProfile) PrepareForValidation(c *handler.Context) *err.ValidationError {
	return err.NewValidationError()
}

func NewUserUpdateProfile(c *handler.Context) (*UserUpdateProfile, error) {
	v := &UserUpdateProfile{}
	e := validation.Body(c, v)
	return v, e
}

// ================================================================
//                       UPDATE EMAIL
// ================================================================

type UserUpdateEmail struct {
	Email string `json:"email"`
}

func (u *UserUpdateEmail) Rules() validation.Rules {
	return validation.Rules{
		"email": {
			"required": {},
			"regex":    {"email"},
			"between":  {"3", "254"},
			"unique":   {"users", "email"},
		},
	}
}

func (u *UserUpdateEmail) PrepareForValidation(c *handler.Context) *err.ValidationError {
	return err.NewValidationError()
}

func NewUserUpdateEmail(c *handler.Context) (*UserUpdateEmail, error) {
	v := &UserUpdateEmail{}
	e := validation.Body(c, v)
	return v, e
}

// ================================================================
//                       UPDATE PASSWORD
// ================================================================

type UserUpdatePassword struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (u *UserUpdatePassword) Rules() validation.Rules {
	return validation.Rules{
		"password": {
			"required": {},
			"between":  {"8", "32"},
		},
		"passwordConfirmation": {
			"required": {},
		},
	}
}

func (u *UserUpdatePassword) PrepareForValidation(c *handler.Context) *err.ValidationError {
	e := err.NewValidationError()
	if u.Password != u.PasswordConfirmation {
		e.AppendM("password", "Las contraseñas no coinciden")
	}
	return e
}

func NewUserUpdatePassword(c *handler.Context) (*UserUpdatePassword, error) {
	v := &UserUpdatePassword{}
	e := validation.Body(c, v)
	return v, e
}

// ================================================================
//                              LOGIN
// ================================================================

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLogin) Rules() validation.Rules {
	return validation.Rules{
		"email": {
			"required": {},
			"regex":    {"email"},
			"between":  {"3", "254"},
		},
		"password": {
			"required": {},
			"between":  {"8", "32"},
		},
	}
}

func (u *UserLogin) PrepareForValidation(c *handler.Context) *err.ValidationError {
	return err.NewValidationError()
}

func NewUserLogin(c *handler.Context) (*UserLogin, error) {
	v := &UserLogin{}
	e := validation.Body(c, v)
	return v, e
}
