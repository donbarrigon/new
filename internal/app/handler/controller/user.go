package controller

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/app/handler/policy"
	"donbarrigon/new/internal/app/handler/service"
	"donbarrigon/new/internal/app/handler/validator"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"

	"golang.org/x/crypto/bcrypt"
)

// post:/api/users
func UserStore(c *handler.Context) {
	dto, e := validator.NewUserStrore(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

	user, e := model.UserCreate(dto)
	if e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory("register", user.ID, user, nil)

	go service.SendEmailConfirm(user)

	if e := auth.SessionStart(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	c.ResponseOk(user)
}

// post:/api/users/login
func UserLogin(c *handler.Context) {
	dto, e := validator.NewUserLogin(c)
	if e != nil {
		c.ResponseError(err.New(err.UNAUTHORIZED, "Credenciales incorrectas", e))
		return
	}

	user, e := model.UserByEmail(dto.Email)
	if e == nil {
		c.ResponseError(err.New(err.UNAUTHORIZED, "Credenciales incorrectas", nil))
		return
	}

	if e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); e != nil {
		c.ResponseError(err.New(err.UNAUTHORIZED, "Credenciales incorrectas", e))
		return
	}

	if e := auth.SessionStart(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory("login", c.Auth.UserID(), user, nil)

	c.ResponseOk(user)
}

// patch:/api/users/profile
func UserUpdateProfile(c *handler.Context) {
	dto, e := validator.NewUserUpdateProfile(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

	user, e := model.UserByHexID(c.Get("id"))
	if e != nil {
		c.ResponseError(e)
		return
	}

	if e := policy.UserUpdate(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	changes, e := user.UpdateProfile(dto)
	if e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory(model.UPDATE_ACTION, c.Auth.UserID(), user, changes)

	c.ResponseOk(user)
}
