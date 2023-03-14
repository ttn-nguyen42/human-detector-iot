package routes

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * POST /api/backend/login
 * Returns a JWT token that will be used to authenticate other requests
 */
func POSTLogin(service services.DeviceInfoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var creds dtos.POSTLoginRequest
		err := ctx.BindJSON(&creds)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: err.Error(),
			})
			return
		}
		if len(creds.DeviceId) == 0 {
			ctx.JSON(http.StatusUnauthorized, MessageResponse{
				Message: "Invalid device_id or password",
			})
			return
		}
		if len(creds.Password) < 8 {
			ctx.JSON(http.StatusUnauthorized, MessageResponse{
				Message: "Invalid device_id or password",
			})
			return
		}
		res, err := service.AuthenticateByPassword(&creds)
		if _, ok := err.(*custom.UnauthorizedError); ok {
			ctx.JSON(http.StatusUnauthorized, MessageResponse{
				Message: "Invalid device_id or password",
			})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, MessageInternalServerError)
		}

		ctx.JSON(http.StatusOK, res)
	}
}

/*
POST /api/backend/register_device
Register a device to the backend
*/
func POSTRegisterDevice(service services.DeviceInfoService, settings services.SettingsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
				Message: "Invalid request body",
			})
			return
		}
		req, err := service.RegisterDevice(&postRequest)
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
		err = settings.CreateSettings(postRequest.DeviceId, &dtos.POSTCreateSettings{})
		if _, ok := err.(*custom.InternalServerError); ok {
			// It is okay to fail, if settings failed to be created
			// it will created next time someone request for it
			ctx.JSON(http.StatusCreated, req)
			return
		}
		ctx.JSON(http.StatusCreated, req)
	}
}
