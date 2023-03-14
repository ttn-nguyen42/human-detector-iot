package routes

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
 * GET /api/backend/settings
 * Requires JWT token. See auths/middleware.go
 */
func GETGetAllSettings(service services.SettingsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		res, err := service.GetSettings(deviceId)
		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		if _, ok := err.(*custom.BadIdError); ok {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: "Device ID not found",
			})
			return
		}
		if _, ok := err.(*custom.ItemNotFoundError); ok {
			// User should already have a settings
			// if not create one with default settings
			err = service.CreateSettings(deviceId, &dtos.POSTCreateSettings{})
		}
		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageResponse{
				Message: "Settings not found, failed to create new default settings",
			})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

/*
 * POST /api/backend/settings/data_rate
 * Requires JWT token. See auths/middleware.go
 */
func POSTUpdateDataRate(service services.SettingsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		var dto dtos.POSTDataRateRequest
		err := ctx.BindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: err.Error(),
			})
			logrus.Info(err.Error())
			return
		}
		err = service.ChangeDataRate(deviceId, dto.RateInSeconds)
		if _, ok := err.(*custom.ItemNotFoundError); ok {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: "Device ID not found",
			})
			return
		}
		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageResponse{
				Message: "Unable to verify device ID",
			})
		}
		if _, ok := err.(*custom.InactiveGatewayError); ok {
			ctx.JSON(http.StatusServiceUnavailable, MessageResponse{
				Message: "Device is not connected with gateway",
			})
			return
		}
		if _, ok := err.(*custom.TimeoutError); ok {
			ctx.JSON(http.StatusGatewayTimeout, MessageResponse{
				Message: "Gateway is not responding",
			})
			return
		}
		if _, ok := err.(*custom.UnableToSendMessage); ok {
			ctx.JSON(http.StatusInternalServerError, MessageResponse{
				Message: "Settings not updated",
			})
			return
		}
		ctx.JSON(http.StatusOK, IdResponse{
			Id: deviceId,
		})
	}
}

/*
 * POST /api/backend/settings/email
 * Requires JWT token
 * */
func POSTChangeNotificationEmail(service services.SettingsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		var dto dtos.POSTChangeEmail
		err := ctx.BindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: err.Error(),
			})
			logrus.Info(err.Error())
			return
		}
		err = service.ChangeEmail(deviceId, dto.Email)
		if _, ok := err.(*custom.FieldMissingError); ok {
			// Service will insert a new "email" field into the database
			// therefore this will not be a matter
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			logrus.Error(err.Error())
			return
		}
		if _, ok := err.(*custom.ItemNotFoundError); ok {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: "Device ID not found",
			})
			return
		}
		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			logrus.Error(err.Error())
			return
		}
		ctx.JSON(http.StatusOK, IdResponse{
			Id: deviceId,
		})
	}
}