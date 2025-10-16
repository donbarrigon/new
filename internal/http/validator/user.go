package validator

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	validate "donbarrigon/new/internal/utils/validation"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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

func (u *UserStore) Rules() validate.Rules {
	return validate.Rules{
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
	return nil
}
