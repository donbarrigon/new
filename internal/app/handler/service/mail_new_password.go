package service

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/config"
)

func SendMailNewPassword(user *model.User, newPassword string) {
	if config.AppName == "es" {
		sendEmailNewPasswordEs(user, newPassword)
	} else {
		sendEmailNewPasswordEn(user, newPassword)
	}
}

func sendEmailNewPasswordEs(user *model.User, newPassword string) {
	subject := "Tu nueva contraseña en " + config.AppName

	body := `
    <h1>Hola ` + user.Profile.Nickname + `</h1>
    <p>Queremos informarte que tu contraseña ha sido restablecida.</p>
    <p>Tu nueva contraseña es:</p>
    <p style="font-size:18px;font-weight:bold;background:#f8f9fa;
              border:1px solid #ddd;padding:10px;border-radius:5px;">
        ` + newPassword + `
    </p>
    <p>Por tu seguridad, te recomendamos cambiarla después de iniciar sesión.</p>
    <p>Si no solicitaste este cambio, por favor contacta con nuestro soporte de inmediato.</p>
    <br>
    <p>Equipo de ` + config.AppName + `</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailNewPasswordEn(user *model.User, newPassword string) {
	subject := "Your new password at " + config.AppName

	body := `
    <h1>Hello ` + user.Profile.Nickname + `</h1>
    <p>We want to let you know that your password has been successfully reset.</p>
    <p>Your new password is:</p>
    <p style="font-size:18px;font-weight:bold;background:#f8f9fa;
              border:1px solid #ddd;padding:10px;border-radius:5px;">
        ` + newPassword + `
    </p>
    <p>For your security, we recommend changing it after logging in.</p>
    <p>If you did not request this change, please contact our support team immediately.</p>
    <br>
    <p>The ` + config.AppName + ` Team</p>
    `

	SendMail(subject, body, user.Email)
}
