package controller

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/app/data/validator"
	"donbarrigon/new/internal/app/handler/policy"
	"donbarrigon/new/internal/app/handler/service"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"

	"golang.org/x/crypto/bcrypt"
)

// get:/api/users
func UserApiIndex(c *handler.Context) {
	if e := policy.UserViewAny(c); e != nil {
		c.ResponseError(e)
		return
	}

	users, e := model.UserPaginate(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

	c.ResponseOk(users)
}

// get:/api/users/show
func UserShow(c *handler.Context) {

	user, e := model.UserByHexID(c.Get("id"))
	if e != nil {
		c.ResponseError(e)
		return
	}

	if e := policy.UserView(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	c.ResponseOk(user)
}

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

	model.CreateHistory("Registro", user.ID, user, nil)

	go service.SendEmailConfirm(user)

	if _, e := auth.SessionStart(c.Writer, c.Request, user); e != nil {
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

	if _, e := auth.SessionStart(c.Writer, c.Request, user); e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory("login", c.Auth.UserID(), user, nil)

	c.ResponseOk(user)
}

// post:/api/users/logout
func UserLogout(c *handler.Context) {
	if e := c.Auth.Destroy(); e != nil {
		c.ResponseError(e)
		return
	}
	c.ResponseNoContent()
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

	model.CreateHistory("Actualización de perfil", c.Auth.UserID(), user, changes)

	c.ResponseOk(user)
}

// ================================================================
//                    ACTUALIZACIÓN DE EMAIL
// ================================================================
// primero se le da al usuario un codigo de verificacion por email
// luego se cambia el email

// post:/api/users/email-code
func UserSendEmailCode(c *handler.Context) {
	service.SendEmailVerificationCode(c.Auth.User.(*model.User), "vc-email")
	model.CreateHistory("Codigo de verificacion cambio de email", c.Auth.UserID(), c.Auth.User.(*model.User), nil)
	c.ResponseNoContent()
}

// patch:/api/users/email
func UserUpdateEmail(c *handler.Context) {
	dto, e := validator.NewUserUpdateEmail(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

	user := c.Auth.User.(*model.User)

	if e := policy.UserUpdate(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	if e := model.TokenExists("vc-email", c.Auth.UserID(), dto.Code); e != nil {
		c.ResponseError(e)
		return
	}

	changes, e := user.UpdateEmail(dto)
	if e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory("Actualización de correo", c.Auth.UserID(), user, changes)

	go service.SendEmailConfirm(user)
	go service.SendEmailChangeRevert(user, changes.Old["email"].(string))

	c.ResponseOk(user)
}

// ================================================================
//                    ACTUALIZACIÓN DE CONTRASEÑA
// ================================================================
// primero se le da al usuario un codigo de verificacion por email
// luego se cambia la password

// post:/api/users/password-code
func UserSendPasswordCode(c *handler.Context) {
	service.SendEmailVerificationCode(c.Auth.User.(*model.User), "vc-password")
	model.CreateHistory("Codigo de verificacion de cambio de contraseña", c.Auth.UserID(), c.Auth.User.(*model.User), nil)
	c.ResponseNoContent()
}

// patch:/api/users/password
func UserUpdatePassword(c *handler.Context) {
	dto, e := validator.NewUserUpdatePassword(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

	user := c.Auth.User.(*model.User)

	if e := policy.UserUpdate(c, user); e != nil {
		c.ResponseError(e)
		return
	}

	if e := model.TokenExists("vc-password", c.Auth.UserID(), dto.Code); e != nil {
		c.ResponseError(e)
		return
	}

	if e := user.UpdatePassword(dto.Password); e != nil {
		c.ResponseError(e)
		return
	}

	model.CreateHistory("Actualización de contraseña", c.Auth.UserID(), user, nil)

	c.ResponseNoContent()
}

// ================================================================
//                  RECUPERACION DE CONTRASEÑA
// ================================================================
// primero hace el forgot password que envia un link al email
// ese link crea una nueva contraseña que se le envia por email al usuario
// luego el usuario cambia la password si quiere

// post:/api/users/forgot-password
func UserForgotPassword(c *handler.Context) {
	service.SendEmailForgotPassword(c.Auth.User.(*model.User))
	model.CreateHistory("Recuperar la contraseña", c.Auth.UserID(), c.Auth.User.(*model.User), nil)
	c.ResponseNoContent()
}

// post:/api/users/new-password
func UserNewPassword(c *handler.Context) {
	user, e := model.UserByHexID(c.Get("u"))
	if e != nil {
		c.ResponseError(e)
		return
	}

	if e := model.TokenExists("forgot-password", user.ID, c.Get("t")); e != nil {
		c.ResponseError(e)
		return
	}

	newPassword, e := user.ResetPassword()
	if e != nil {
		c.ResponseError(e)
		return
	}

	go service.SendMailNewPassword(user, newPassword)

	model.CreateHistory("Restablecimiento de contraseña", c.Auth.UserID(), user, nil)

	c.ResponseNoContent()

}
