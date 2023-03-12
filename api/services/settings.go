package services

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/models"
	"iot_api/repositories"

	"github.com/jinzhu/copier"
)

type SettingsService interface {
	CreateSettings(deviceId string, req *dtos.POSTCreateSettings) error
	GetSettings(deviceId string) (*dtos.GETGetSettings, error)
	ChangeDataRate(deviceId string, newRate int) error
}

type settingsService struct {
	repo      repositories.SettingsRepository
	credsRepo repositories.DeviceInfoRepository
	command   CommandService
}

func NewSettingsService(repo repositories.SettingsRepository,
	credsRepo repositories.DeviceInfoRepository,
	commandService CommandService) SettingsService {
	return &settingsService{
		repo:      repo,
		credsRepo: credsRepo,
		command:   commandService,
	}
}

func (s *settingsService) CreateSettings(deviceId string, req *dtos.POSTCreateSettings) error {
	_, err := s.credsRepo.GetCredentials(deviceId)
	if err != nil {
		return custom.NewInternalServerError(err.Error())
	}
	if _, ok := err.(*custom.ItemNotFoundError); ok {
		return custom.NewItemNotFoundError("Device not found")
	}
	mod := &models.Settings{
		DeviceId: deviceId,
		DataRate: req.DataRateInSeconds,
	}
	_, err = s.repo.SaveSettings(mod)
	if err != nil {
		return custom.NewInternalServerError("Unable to save settings")
	}
	return nil
}

func (s *settingsService) GetSettings(deviceId string) (*dtos.GETGetSettings, error) {
	_, err := s.credsRepo.GetCredentials(deviceId)
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	if _, ok := err.(*custom.ItemNotFoundError); ok {
		return nil, custom.NewBadIdError("Device not found")
	}
	mod, err := s.repo.GetSettings(deviceId)
	var dto dtos.GETGetSettings
	copier.Copy(&dto, &mod)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (s *settingsService) ChangeDataRate(deviceId string, newRate int) error {
	_, err := s.credsRepo.GetCredentials(deviceId)
	if err != nil {
		return custom.NewInternalServerError(err.Error())
	}
	if _, ok := err.(*custom.ItemNotFoundError); ok {
		return custom.NewItemNotFoundError("Device not found")
	}
	// Send command first
	// if device received request, update in the database
	err = s.command.SendDataRateUpdate(deviceId, dtos.POSTDataRateRequest{
		RateInSeconds: newRate,
	})
	if err != nil {
		return err
	}
	// Key equals bson key in the Settings model for data rate
	err = s.repo.UpdateSettings(deviceId, "data_rate", newRate)
	if err != nil {
		return err
	}
	return nil
}
