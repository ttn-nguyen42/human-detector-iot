package repositories

import (
	"iot_api/database"
	"iot_api/models"
)

type DeviceInfoRepository interface {

}

type deviceInfoRepository struct {
	col database.Collection[models.DeviceCredentials]
}

func NewDeviceInfoRepository(db database.Collection[models.DeviceCredentials]) DeviceInfoRepository {
	return &deviceInfoRepository{
		col: db,
	}
}