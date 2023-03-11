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
GET /api/backend/check_active
Check if the provided device is active or not
*/
func GetIsDeviceActive(service services.CommandService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Device ID from middleware
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		err := service.SendStatusCheck(deviceId)
		if _, ok := err.(*custom.InactiveGatewayError); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		if _, ok := err.(*custom.TimeoutError); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		if _, ok := err.(*custom.UnableToSendMessage); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		ctx.JSON(http.StatusOK, IdResponse{
			Id: deviceId,
		})
	}
}

/*
POST /api/backend/settings/data_rate
Requires JWT token. See auths/middleware.go
*/
func POSTUpdateDataRate(service services.CommandService) gin.HandlerFunc {
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
		err = service.SendDataRateUpdate(deviceId, dto)
		if _, ok := err.(*custom.InactiveGatewayError); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		if _, ok := err.(*custom.TimeoutError); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		if _, ok := err.(*custom.UnableToSendMessage); ok {
			ctx.JSON(http.StatusServiceUnavailable, IdResponse{
				Id: deviceId,
			})
			return
		}
		// Unimplemented
		ctx.JSON(http.StatusOK, MessageResponse{
			Message: "Ok",
		})
	}
}
