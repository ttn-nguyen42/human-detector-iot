package routes

import (
	"iot_api/custom"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
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
