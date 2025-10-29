package service

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
)

func SendEmailConfirm(user *model.User) {
	tk, e := model.TokenCreate(user.ID, "confirm-email", map[string]string{})
	if e != nil {
		logs.Error("Error al crear el token para confirmar el correo: %s", e.Error())
		return
	}
	url := config.AppURL + "/users/confirm-email?u=" + user.ID.Hex() + "&t=" + tk.Token
	switch config.AppLocale {
	case "es":
		sendEmailConfirmEs(user, url)
	default:
		sendEmailConfirmEn(user, url)
	}
}

func sendEmailConfirmEs(user *model.User, url string) {

	subject := "Confirma tu cuenta en " + config.AppName

	body := `
    <h1>Bienvenido a ` + config.AppName + `</h1>
    <p>Gracias por registrarte. Para completar tu registro, haz clic en el siguiente enlace:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirmar mi correo
        </a>
    </p>
    <p>Si no fuiste tú quien se registró, puedes ignorar este mensaje.</p>
	<p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + url + `</p>
    <br>
    <p>Equipo de ` + config.AppName + `</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailConfirmEn(user *model.User, url string) {

	subject := "Confirm your account on " + config.AppName

	body := `
    <h1>Welcome to ` + config.AppName + `</h1>
    <p>Thank you for signing up. To complete your registration, please click the link below:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirm my email
        </a>
    </p>
    <p>If you did not create this account, you can safely ignore this message.</p>
	<p>If you cannot click the button, copy and paste this link into your browser:</p>
    <p>` + url + `</p>
    <br>
    <p>The ` + config.AppName + ` Team</p>
    `

	SendMail(subject, body, user.Email)
}
