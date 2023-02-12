package routes

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*

 */
func POSTRegisterDevice(service services.DeviceInfoService) gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		/*
		POST /api/backend/register_device
		
		Handles device registration.
		The device will send its device_id
		and the backend returns a password associated with it
		*/
		var postRequest dtos.POSTRegisterDeviceDto
		// Take in raw payload in bytes
		// deserializes it into a normal struct
		err := ctx.BindJSON(&postRequest)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: err.Error(),
			})
		}

		req, err := service.CreatePassword(&postRequest)

		if _, ok := err.(*custom.AlreadyRegisteredError); ok {
			ctx.JSON(http.StatusConflict, MessageResponse{
				Message: "Already registered device",
			})
			return
		}

		if _, ok := err.(*custom.InternalServerError); ok {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, req)
	}
}