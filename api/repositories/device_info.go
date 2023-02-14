package repositories

import (
	"iot_api/custom"
	"iot_api/database"
	"iot_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeviceInfoRepository interface {
	SaveCredentials(ent *models.DeviceCredentials) (string, error)
	GetCredentials(deviceId string) (*models.DeviceCredentials, error)
}

type deviceInfoRepository struct {
	col database.Collection[models.DeviceCredentials]
}

func NewDeviceInfoRepository(db database.Collection[models.DeviceCredentials]) DeviceInfoRepository {
	return &deviceInfoRepository{
		col: db,
	}
}

func (r *deviceInfoRepository) SaveCredentials(ent *models.DeviceCredentials) (string, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	resourceId, err := r.col.InsertOne(ctx, *ent)
	if err != nil {
		return "", custom.NewInternalServerError(err.Error())
	}
	return resourceId, err
}

func (r *deviceInfoRepository) GetCredentials(deviceId string) (*models.DeviceCredentials, error) {
	ctx, cancel := database.GetContext()
	defer cancel()
	filter := bson.D{
		{Key: "device_id", Value: deviceId},
	}
	result := &models.DeviceCredentials{}
	err := r.col.FindOne(ctx, result, filter)
	if err == mongo.ErrNoDocuments {
		return nil, custom.NewItemNotFoundError(err.Error())
	}
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	return result, nil
}
