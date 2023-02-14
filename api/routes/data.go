package routes

import (
	"fmt"
	"io"
	"iot_api/models"
	"iot_api/network"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		topic := fmt.Sprintf("%v/%v", "yolobit/data/sensor", deviceId)
		res := make(chan interface{})
		go func() {
			network.GetClient().Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
				// Is a duplicate
				if m.Duplicate() {
					return
				}
				var parsedPayload models.SensorData
				err := models.StructifyBytes(m.Payload(), &parsedPayload)
				if err != nil {
					return
				}
				res <- parsedPayload
			})
		}()
		ended := ctx.Stream(func(w io.Writer) bool {
			if data, ok := <-res; ok {
				ctx.SSEvent("data", data)
				return true
			}
			return false
		})
		if ended {
			close(res)
		}
	}
}

/*
POST /api/backend/settings/data_rate
Requires JWT token. See auths/middleware.go
*/
func POSTUpdateDataRate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
GET /api/backend/settings
Requires JWT token. See auths/middleware.go
*/
func GETGetAllSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
