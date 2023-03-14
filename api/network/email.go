package network

import (
	"iot_api/utils"
	"net/smtp"
	"sync"

	"github.com/sirupsen/logrus"
)

var emailOnce sync.Once

type SMTPAuth struct {
	Auth smtp.Auth
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
		smtpAuth = &SMTPAuth{
			Auth: smtp.PlainAuth("", creds.Username, creds.Password, creds.Host),
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
