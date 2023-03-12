package routes

import (
	"fmt"
	"io"
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/models"
	"iot_api/network"
	"iot_api/services"
	"iot_api/utils"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
GET /api/backend/settings
Requires JWT token. See auths/middleware.go
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
POST /api/backend/settings/data_rate
Requires JWT token. See auths/middleware.go
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
