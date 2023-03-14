package routes

import (
	"io"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * GET /api/backend/data
 * Requires JWT token. See auths/middleware.go
*/
func GETGetDeviceData(service services.DataService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Device ID from middleware
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		quit := make(chan bool)
		stopStream, res := service.RetrieveSensorDataStream(deviceId, quit)
		ended := ctx.Stream(func(w io.Writer) bool {
			if data, ok := <-res; ok {
				ctx.SSEvent("data", data)
				return true
			}
			return false
		})
		if ended {
			quit <- true
			stopStream()
			close(quit)
		}
	}
}