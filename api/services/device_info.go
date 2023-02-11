package services

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/models"
	"iot_api/repositories"
	"iot_api/utils"
)

type DeviceInfoService interface {
	CreatePassword(req *dtos.POSTRegisterDeviceDto) (*dtos.POSTRegisterDeviceResponse, error)
}

type deviceInfoService struct {
	repo repositories.DeviceInfoRepository
}

func NewDeviceInfoService(repo repositories.DeviceInfoRepository) DeviceInfoService {
	return &deviceInfoService{
		repo: repo,
	}
}

func (s *deviceInfoService) CreatePassword(req *dtos.POSTRegisterDeviceDto) (*dtos.POSTRegisterDeviceResponse, error) {
	_, err := s.repo.GetCredentials(req.DeviceId)
	if err == nil {
		return nil, custom.NewAlreadyRegisteredError("Device already registered")
	}
	if _, ok := err.(*custom.ItemNotFoundError); !ok {
		return nil, custom.NewInternalServerError(err.Error())
	}
	rawPass := utils.GetRandomUUID()
	hashPass := utils.MD5Hash(rawPass)
	ent := &models.DeviceCredentials{
		DeviceId: req.DeviceId,
		Password: hashPass,
	}
	_, err = s.repo.SaveCredentials(ent)
	if err != nil {
		return &dtos.POSTRegisterDeviceResponse{
			Password: rawPass,
		}, err
	}
	return &dtos.POSTRegisterDeviceResponse{
		Password: rawPass,
	}, nil
}
