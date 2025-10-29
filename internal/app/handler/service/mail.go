package service

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
	"net/smtp"
	"strings"
	"time"
)

func SendMail(subject string, body string, to ...string) {

	if config.MailUsername == "tuemail@gmail.com" {
		logs.Warning("Failed to send email: no email configured")
		return
	}

	for i := range 3 {

		auth := smtp.PlainAuth(config.MailIdentity, config.MailUsername, config.MailPassword, config.MailHost)

		msg := []byte(
			"From: " + config.MailFromName + " <" + config.MailUsername + ">\r\n" +
				"To: " + strings.Join(to, ",") + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"MIME-Version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		e := smtp.SendMail(config.MailHost+":"+config.MailPort, auth, config.MailUsername, to, msg)
		if e != nil {
			logs.Error("Failed to send email: try: %d to: %s error: %s", i, to, e.Error())
			time.Sleep(30 * time.Second)
			continue
		}

		return
	}
}
