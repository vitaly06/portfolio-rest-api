package mailer

import (
	"fmt"
	"net/smtp"

	"github.com/vitaly06/portfolio-rest-api/internal/config"
)

type Mailer struct {
	cfg *config.Config
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendContactEmails(userName, userEmail, userPhone, comment, aiReply string) error {
	auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPass, m.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", m.cfg.SMTPHost, m.cfg.SMTPPort)

	// 1. Письмо владельцу (Виталию)
	ownerSubject := "Subject: Новая заявка с сайта-портфолио!\r\n"
	ownerBody := fmt.Sprintf(
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n%s"+
			"<h3>Новое обращение от %s</h3>\n"+
			"<p><b>Email:</b> %s</p>\n"+
			"<p><b>Телефон:</b> %s</p>\n"+
			"<p><b>Комментарий:</b> %s</p>\n",
		ownerSubject, userName, userEmail, userPhone, comment,
	)

	err := smtp.SendMail(addr, auth, m.cfg.SMTPUser, []string{m.cfg.OwnerEmail}, []byte(ownerBody))
	if err != nil {
		return fmt.Errorf("ошибка отправки владельцу: %w", err)
	}

	// 2. Письмо пользователю (Автоответ)
	userSubject := "Subject: Ваше обращение принято — Виталий (Golang Developer)\r\n"
	userBody := fmt.Sprintf(
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n%s"+
			"Приветствую, %s!<br><br>\n"+
			"%s<br><br>\n"+
			"---<br>\n"+
			"С уважением, Виталий.<br>\n"+
			"Моё портфолио и контакты на hh.ru.",
		userSubject, userName, aiReply,
	)

	return smtp.SendMail(addr, auth, m.cfg.SMTPUser, []string{userEmail}, []byte(userBody))
}
