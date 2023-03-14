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
 * A webhook for sending notification email
 * POST /api/backend/notify
 * Requires JWT authentication
 */
func POSTSendEmailNotification(service services.NotificationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetString("device_id")
		if len(deviceId) == 0 {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}
		var dto dtos.POSTSendEmailNotification
		err := ctx.BindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: err.Error(),
			})
			return
		}
		err = service.SendEmailNotification(deviceId, dto.Title, dto.Message)
		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			logrus.Error(err.Error())
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			logrus.Error(err.Error())
			return
		}
		ctx.JSON(http.StatusOK, IdResponse{
			Id: deviceId,
		})
	}
}
