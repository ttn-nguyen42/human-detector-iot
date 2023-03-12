package repositories

import (
	"iot_api/custom"
	"iot_api/database"
	"iot_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SettingsRepository interface {
	SaveSettings(ent *models.Settings) (string, error)
	GetSettings(deviceId string) (*models.Settings, error)
	UpdateSettings(deviceId string, key string, value interface{}) error
}

type settingsRepository struct {
	col database.Collection[models.Settings]
}

func NewSettingsRepository(db database.Collection[models.Settings]) SettingsRepository {
	return &settingsRepository{
		col: db,
	}
}

func (r *settingsRepository) SaveSettings(ent *models.Settings) (string, error) {
	ctx, cancel := database.GetContext()
	defer cancel()
	resourceId, err := r.col.InsertOne(ctx, *ent)
	if err != nil {
		return "", custom.NewInternalServerError(err.Error())
	}
	return resourceId, err
}

func (r *settingsRepository) GetSettings(deviceId string) (*models.Settings, error) {
	ctx, cancel := database.GetContext()
	defer cancel()
	filter := bson.D{
		{Key: "device_id", Value: deviceId},
	}
	result := &models.Settings{}
	err := r.col.FindOne(ctx, result, filter)
	if err == mongo.ErrNoDocuments {
		return nil, custom.NewItemNotFoundError(err.Error())
	}
	if err != nil {
		return nil, custom.NewInternalServerError(err.Error())
	}
	return result, nil
}

func (r *settingsRepository) UpdateSettings(deviceId string, key string, value interface{}) error {
	ctx, cancel := database.GetContext()
	defer cancel()
	filter := bson.D{{Key: "device_id", Value: deviceId}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   key,
			Value: value,
		}}},
	}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil || res == nil {
		return custom.NewInternalServerError(err.Error())
	}
	if res.MatchedField == 0 {
		return custom.NewFieldMissingError("Field not found")
	}
	if res.MatchedFilter == 0 {
		return custom.NewItemNotFoundError("Cannot find filtered document")
	}
	return nil
}
