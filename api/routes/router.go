package routes

import (
	"iot_api/database"
	"iot_api/models"
	"iot_api/repositories"
	"iot_api/services"

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

	deviceInfoService := services.NewDeviceInfoService(deviceInfoRepo)

	v1 := engine.Group("/api/backend")

	v1.POST("/register_device", POSTRegisterDevice(deviceInfoService))

}