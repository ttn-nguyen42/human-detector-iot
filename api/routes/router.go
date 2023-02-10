package routes

import (
	"iot_api/database"
	"iot_api/models"
	repositories "iot_api/repositories/device_info"
	deviceInfo "iot_api/services/device_info"

	"github.com/gin-gonic/gin"
)

const (
	DatabaseName = "iot_database"
	DeviceInfoCollection = "device_info"
)

func Create(engine *gin.Engine) {
	// Enables logging for HTTP server
	engine.Use(gin.Logger())
	// Returns 500 on panic()
	engine.Use(gin.Recovery())

	dbClient := database.GetClient()

	// Device Info collection
	deviceInfoCol := &database.MongoCollection[models.DeviceCredentials]{
		Col: dbClient.Database(DatabaseName).Collection(DeviceInfoCollection),
	}
	
	deviceInfoRepo := repositories.NewDeviceInfoRepository(deviceInfoCol)

	deviceInfoService := deviceInfo.NewDeviceInfoService(deviceInfoRepo)

	

}