package service

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
)

func SendEmailVerificationCode(user *model.User, action string) {
	tk, e := model.TokenCreateVerificationCode(user.ID, action, map[string]string{})
	if e != nil {
		logs.Error("Error al crear el token para confirmar el correo: %s", e.Error())
		return
	}
	url := config.AppURL + "/users/" + action + "?u=" + user.ID.Hex() + "&t=" + tk.Token
	switch config.AppLocale {
	case "es":
		sendEmailVerificationCodeEs(user, url)
	default:
		sendEmailVerificationCodeEn(user, url)
	}
}

func sendEmailVerificationCodeEs(user *model.User, code string) {

	subject := "Código de verificación - " + config.AppName

	body := `
    <h1>Código de verificación</h1>
    <p>Hola, has solicitado un código de verificación para tu cuenta en ` + config.AppName + `.</p>
    <p>Tu código de verificación es:</p>
    <div style="text-align:center;margin:30px 0;">
        <span style="display:inline-block;padding:15px 30px;background:#f0f0f0;
                     font-size:32px;font-weight:bold;letter-spacing:8px;
                     border-radius:8px;color:#333;">
           ` + code + `
        </span>
    </div>
    <p><strong>Este código expirará en 10 minutos.</strong></p>
    <p>Si no solicitaste este código, puedes ignorar este mensaje de forma segura.</p>
    <br>
    <p>Equipo de ` + config.AppName + `</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailVerificationCodeEn(user *model.User, code string) {

	subject := "Verification Code - " + config.AppName

	body := `
    <h1>Verification Code</h1>
    <p>Hello, you have requested a verification code for your account on ` + config.AppName + `.</p>
    <p>Your verification code is:</p>
    <div style="text-align:center;margin:30px 0;">
        <span style="display:inline-block;padding:15px 30px;background:#f0f0f0;
                     font-size:32px;font-weight:bold;letter-spacing:8px;
                     border-radius:8px;color:#333;">
           ` + code + `
        </span>
    </div>
    <p><strong>This code will expire in 10 minutes.</strong></p>
    <p>If you did not request this code, you can safely ignore this message.</p>
    <br>
    <p>The ` + config.AppName + ` Team</p>
    `

	SendMail(subject, body, user.Email)
}
