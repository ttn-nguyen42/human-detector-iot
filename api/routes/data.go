package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
GET /api/backend/data
Requires JWT token. See auths/middleware.go
*/
func GETGetDeviceData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Device ID from middleware
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		// Unimplemented
		ctx.JSON(http.StatusOK, MessageResponse{
			Message: "Ok",
		})
	}
}

/*
POST /api/backend/settings/data_rate
Requires JWT token. See auths/middleware.go
*/
func POSTUpdateDataRate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0{
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		// Unimplemented
		ctx.JSON(http.StatusOK, MessageResponse{
			Message: "Ok",
		})
	}
}