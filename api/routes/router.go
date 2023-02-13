package routes

import (
	"iot_api/auths"
	"iot_api/database"
	"iot_api/models"
	"iot_api/network"
	"iot_api/repositories"
	"iot_api/services"

	"github.com/gin-gonic/gin"
)

const (
	DatabaseName         = "iot_database"
	DeviceInfoCollection = "device_info"
)

func Create(engine *gin.Engine) {
	// Enables logging for HTTP server
	engine.Use(gin.Logger())
	// Returns 500 on panic()
	engine.Use(gin.Recovery())

	dbClient := database.GetClient()
	network.GetClient()

	// Device Info collection
	deviceInfoCol := &database.MongoCollection[models.DeviceCredentials]{
		Col: dbClient.Database(DatabaseName).Collection(DeviceInfoCollection),
	}

	deviceInfoRepo := repositories.NewDeviceInfoRepository(deviceInfoCol)

	deviceInfoService := services.NewDeviceInfoService(deviceInfoRepo)

	// Unprotected endpoints
	public := engine.Group("/api/backend")
	public.POST("/register_device", POSTRegisterDevice(deviceInfoService))
	public.POST("/login", POSTLogin(deviceInfoService))

	// Protected endpoints
	// Using JWTs
	// See auths/middleware
	protected := engine.Group("/api/backend")
	protected.Use(auths.JwtAuthMiddleware())
	protected.GET("/data", ssEventHeader(), GETGetDeviceData())

	protected.POST("/settings/data_rate", POSTUpdateDataRate())
	protected.GET("/settings", GETGetAllSettings())
}
