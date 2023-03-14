package routes

import (
	"iot_api/auths"
	"iot_api/database"
	"iot_api/models"
	"iot_api/network"
	"iot_api/repositories"
	"iot_api/services"
	"net/http"

	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/gin-gonic/gin"
)

const (
	DatabaseName         = "iot_database"
	DeviceInfoCollection = "device_info"
	SettingsCollection   = "settings"
)

func Create(engine *gin.Engine) {
	// Enables logging for HTTP server
	engine.Use(gin.Logger())
	// Returns 500 on panic()
	engine.Use(gin.Recovery())

	// Healthcheck
	engine.Use(healthcheck.New(healthcheck.Config{
		HeaderName:   "X-Check",
		HeaderValue:  "healthcheck",
		ResponseCode: http.StatusTeapot,
		ResponseText: "im a teapot",
	}))

	dbClient := database.GetClient()
	network.GetClient()
	smtpAuth := network.GetSMTPAuth()

	// Device Info collection
	deviceInfoCol := &database.MongoCollection[models.DeviceCredentials]{
		Col: dbClient.Database(DatabaseName).Collection(DeviceInfoCollection),
	}

	settingsCol := &database.MongoCollection[models.Settings]{
		Col: dbClient.Database(DatabaseName).Collection(SettingsCollection),
	}

	deviceInfoRepo := repositories.NewDeviceInfoRepository(deviceInfoCol)
	settingsRepo := repositories.NewSettingsRepository(settingsCol)

	deviceInfoService := services.NewDeviceInfoService(deviceInfoRepo)
	commandService := services.NewCommandService("yolobit/command/activity", "yolobit/command/response")
	settingsService := services.NewSettingsService(settingsRepo, deviceInfoRepo, commandService)
	dataService := services.NewDataService()

	notificationService := services.NewNotificationService(smtpAuth, settingsService)

	/*
	 * Unprotected endpoints
	 */
	public := engine.Group("/api/backend")
	public.POST("/register_device", POSTRegisterDevice(deviceInfoService, settingsService))
	public.POST("/login", POSTLogin(deviceInfoService))

	/*
	 * Protected endpoints
	 */
	protected := engine.Group("/api/backend")
	protected.Use(auths.JwtAuthMiddleware())

	protected.POST("/settings/data_rate", POSTUpdateDataRate(settingsService)) // Only works when device is connected
	protected.POST("/settings/email", POSTChangeNotificationEmail(settingsService))
	protected.GET("/settings", GETGetAllSettings(settingsService))

	protected.GET("/data", ssEventHeader(), GETGetDeviceData(dataService))
	protected.GET("/check_active", GetIsDeviceActive(commandService))
	protected.POST("/notify", POSTSendEmailNotification(notificationService))
}
