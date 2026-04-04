package apiv0

import (
	"fmt"
	"net/smtp"
)

const (
	emailFrom = "no-reply@dublinbikeparking.com"
	emailTo   = "hello@katiechapman.ie"
)

func sendMail(cfg *smtpConfig, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	msg := fmt.Sprintf(
		"From: DublinBikeParking <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		emailFrom, emailTo, subject, body,
	)
	return smtp.SendMail(addr, auth, emailFrom, []string{emailTo}, []byte(msg))
}
