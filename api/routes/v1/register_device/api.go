package register_device

import (
	deviceInfo "iot_api/services/device_info"

	"github.com/gin-gonic/gin"
)

func POSTRegisterDevice() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		/*
		POST /api/v1/register_device
		
		Handles device registration.
		The device will send its device_id
		and the backend returns a password associated with it
		*/
		var postRequest deviceInfo.POSTRegisterDeviceDto
		// Take in raw payload in bytes
		// deserializes it into a normal struct
		ctx.BindJSON(&postRequest)
	}
}