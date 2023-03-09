package routes

import (
	"fmt"
	"io"
	"iot_api/models"
	"iot_api/network"
	"iot_api/utils"
	"net/http"
	"time"

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
		quit := make(chan bool)
		if !utils.IsTestMode() {
			/* If in production mode, subscribes to MQTT server */
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
		} else {
			/* In test mode, data is generated */
			go func() {
				for {
					select {
					case <-quit:
						return
					default:
						res <- utils.GetRandomSensorData(deviceId)
						time.Sleep(time.Second * 2) // Send new data every 2 seconds
					}
				}
			}()
		}

		ended := ctx.Stream(func(w io.Writer) bool {
			if data, ok := <-res; ok {
				ctx.SSEvent("data", data)
				return true
			}
			return false
		})
		if ended {
			network.GetClient().Unsubscribe(topic)
			quit <- true
			close(quit)
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
