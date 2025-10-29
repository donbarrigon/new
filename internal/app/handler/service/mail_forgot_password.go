package service

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
)

func SendEmailForgotPassword(user *model.User) {
	tk, e := model.TokenCreate(user.ID, "forgot-password", map[string]string{})
	if e != nil {
		logs.Error("Error al crear el token para restablecer la contraseña: %s", e.Error())
		return
	}
	url := config.AppURL + "/users/forgot-password?u=" + user.ID.Hex() + "&t=" + tk.Token

	if config.AppLocale == "es" {
		sendEmailForgotPasswordEs(user, url)
	} else {
		sendEmailForgotPasswordEn(user, url)
	}
}

func sendEmailForgotPasswordEs(user *model.User, url string) {

	subject := "Restablece tu contraseña en " + config.AppName

	body := `
    <h1>Hola ` + user.Profile.Nickname + `</h1>
    <p>Recibimos una solicitud para restablecer tu contraseña en ` + config.AppName + `.</p>
    <p>Se creara una nueva contraseña haciendo clic en el siguiente enlace:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Restablecer contraseña
        </a>
    </p>
    <p>Si no solicitaste este cambio, por favor restablece tu contraseña inmediatamente o contacta a nuestro soporte.</p>
    <p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + url + `</p>
    <br>
    <p>Equipo de ` + config.AppName + `</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailForgotPasswordEn(user *model.User, url string) {

	subject := "Reset your password at " + config.AppName

	body := `
    <h1>Hello ` + user.Profile.Nickname + `</h1>
    <p>We received a request to reset your password at ` + config.AppName + `.</p>
    <p>A new password will be created by clicking the link below:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Reset Password
        </a>
    </p>
    <p>If you did not request this change, make this change, please reset your password immediately or contact our support team.</p>
    <p>If you cannot click the button, copy and paste this link into your browser:</p>
    <p>` + url + `</p>
    <br>
    <p>The ` + config.AppName + ` Team</p>
    `

	SendMail(subject, body, user.Email)
}
