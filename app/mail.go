package app

import (
	"fmt"
	"os"
	"time"

	gomail "gopkg.in/mail.v2"
)

func (a App) SendLog(logfile string, stopped bool) error {
	text := ""
	if stopped {
		text = "Es gab Fehler!"
	} else {
		text = "Keine Fehler"
	}
	m := gomail.NewMessage()
	m.SetHeader("From", a.env.LOG_MAIL)
	m.SetHeader("To", a.env.LOG_MAIL)
	m.SetHeader("Subject", fmt.Sprintf("Log vom %s", time.Now().Local().String()))
	m.Attach("log.txt")
	m.SetBody("text/html", text)

	d := gomail.NewDialer(a.env.MAIL_SERVER, a.env.MAIL_PORT, a.env.MAIL_USER, a.env.MAIL_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	os.Remove(logfile)

	return nil
}
