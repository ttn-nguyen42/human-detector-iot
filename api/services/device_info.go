package services

import (
	"iot_api/auths"
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/models"
	"iot_api/repositories"
	"iot_api/utils"

	"github.com/sirupsen/logrus"
)

type DeviceInfoService interface {
	CreatePassword(req *dtos.POSTRegisterDeviceDto) (*dtos.POSTRegisterDeviceResponse, error)
	AuthenticateByPassword(req *dtos.POSTLoginRequest) (*dtos.POSTLoginResponse, error)
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
	hashPass := utils.GetPasswordHash(rawPass)
	ent := &models.DeviceCredentials{
		DeviceId: req.DeviceId,
		Password: hashPass,
	}
	_, err = s.repo.SaveCredentials(ent)
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	return &dtos.POSTRegisterDeviceResponse{
		Password: rawPass,
	}, nil
}

func (s *deviceInfoService) AuthenticateByPassword(req *dtos.POSTLoginRequest) (*dtos.POSTLoginResponse, error) {
	creds, err := s.repo.GetCredentials(req.DeviceId)
	if _, ok := err.(*custom.ItemNotFoundError); ok {
		return nil, custom.NewUnauthorizedError("User not found")
	}
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	logrus.Debug(creds.Password)
	hashedInput := utils.GetPasswordHash(req.Password)
	logrus.Debug(hashedInput)
	if hashedInput != creds.Password {
		return nil, custom.NewUnauthorizedError("Wrong password")
	}
	token, err := auths.GenerateJwt(req.DeviceId)
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	return &dtos.POSTLoginResponse{
		Token: token,
	}, nil
}
