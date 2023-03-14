package network

import (
	"iot_api/utils"
	"strconv"
	"sync"

	"github.com/go-mail/mail"
	"github.com/sirupsen/logrus"
)

var emailOnce sync.Once

type SMTPAuth struct {
	Dialer *mail.Dialer
	Credentials *utils.SMTPSettings
}

// Singleton
var smtpAuth *SMTPAuth

func GetSMTPAuth() *SMTPAuth {
	if smtpAuth != nil {
		return smtpAuth
	}
	emailOnce.Do(func() {
		creds := getSMTPSettings()
		port, _ := strconv.Atoi(creds.Port)
		smtpAuth = &SMTPAuth{
			Dialer: mail.NewDialer(creds.Host, port, creds.Username, creds.Password),
			Credentials: creds,
		}
	})
	return smtpAuth
}

func getSMTPSettings() *utils.SMTPSettings {
	sets, err := utils.GetSMTPSettings()
	if err != nil {
		logrus.Fatalf("Unable to get SMTP authentications: %v", err.Error())
		return nil
	}
	return sets
}
