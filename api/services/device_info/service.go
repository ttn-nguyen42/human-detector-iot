package deviceInfo

import repositories "iot_api/repositories/device_info"

type DeviceInfoService interface {

}

type deviceInfoService struct {
	repo repositories.DeviceInfoRepository
}

func NewDeviceInfoService(repo repositories.DeviceInfoRepository) DeviceInfoService {
	return &deviceInfoService{
		repo: repo,
	}
}