package service

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
)

func SendEmailChangeRevert(user *model.User, oldEmail string) {
	tk, e := model.TokenCreate(user.ID, "email-change-revert", map[string]string{"old_email": oldEmail})
	if e != nil {
		logs.Error("Error al crear el token para revertir el cambio de correo: %s", e.Error())
		return
	}
	url := config.AppURL + "/users/email-change-revert?u=" + user.ID.Hex() + "&t=" + tk.Token

	switch config.AppLocale {
	case "es":
		sendEmailChangeRevertEs(user, oldEmail, url)
	default:
		sendEmailChangeRevertEn(user, oldEmail, url)
	}
}

func sendEmailChangeRevertEs(user *model.User, oldEmail string, url string) {

	subject := "Tu correo en " + config.AppName + " ha sido actualizado"
	body := `
    <h1>Hola de nuevo en ` + config.AppName + `</h1>
    <p>Queremos informarte que tu dirección de correo fue actualizada recientemente.</p>
    <p>Nuevo correo: <strong>` + user.Email + `</strong></p>
    <p>Si realizaste este cambio, no necesitas hacer nada.</p>
    <p>Pero si <strong>NO fuiste tú</strong>, puedes revertir el cambio haciendo clic en el siguiente enlace:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revertir cambio de correo
        </a>
    </p>
    <p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + url + `</p>
    <br>
    <p>Equipo de ` + config.AppName + `</p>
    `

	// Se envía al email ANTIGUO, no al nuevo
	SendMail(subject, body, oldEmail)
}

func sendEmailChangeRevertEn(user *model.User, oldEmail string, url string) {

	subject := "Your email address on " + config.AppName + " has been updated"

	body := `
    <h1>Hello again from ` + config.AppName + `</h1>
    <p>We want to let you know that your email address was recently updated.</p>
    <p>New email: <strong>` + user.Email + `</strong></p>
    <p>If you made this change, no further action is required.</p>
    <p>But if <strong>you did NOT make this change</strong>, you can revert it by clicking the link below:</p>
    <p>
        <a href="` + url + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revert email change
        </a>
    </p>
    <p>If the button doesn’t work, copy and paste this link into your browser:</p>
    <p>` + url + `</p>
    <br>
    <p>The ` + config.AppName + ` Team</p>
    `

	// Send to the OLD email, not the new one
	SendMail(subject, body, oldEmail)
}
