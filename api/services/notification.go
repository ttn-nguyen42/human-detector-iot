package services

import (
	"fmt"
	"iot_api/custom"
	"iot_api/network"
	"net/smtp"
)

type SendEmailOptions struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
}

type NotificationService interface {
	SendEmail(dto *SendEmailOptions) error
	SendEmailNotification(deviceId string, title string, message string) error
}

type notificationService struct {
	Credentials     *network.SMTPAuth
	SettingsService SettingsService
}

func NewNotificationService(creds *network.SMTPAuth, settingsService SettingsService) NotificationService {
	return &notificationService{
		Credentials:     creds,
		SettingsService: settingsService,
	}
}

func (s *notificationService) SendEmail(dto *SendEmailOptions) error {
	creds := s.Credentials.Credentials
	addr := fmt.Sprintf("%v:%v", creds.Host, creds.Port)
	sender := s.Credentials.Credentials.Sender
	if len(dto.From) > 0 {
		sender = dto.From
	}
	err := smtp.SendMail(addr, s.Credentials.Auth, sender, dto.To, []byte(dto.Message))
	if err != nil {
		return custom.NewInternalServerError(err.Error())
	}
	return nil
}

func (s *notificationService) SendEmailNotification(deviceId string, title string, message string) error {
	settings, err := s.SettingsService.GetSettings(deviceId)
	if err != nil {
		return err
	}
	receivers := settings.NotificationEmails
	err = s.SendEmail(&SendEmailOptions{
		To:      receivers,
		Title:   title,
		Message: message,
	})
	if err != nil {
		return custom.NewInternalServerError(err.Error())
	}
	return nil
}
