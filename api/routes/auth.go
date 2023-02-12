package routes

import (
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
For web app
POST /api/backend/login
Body {
	"device_id": "something",
	"password": "rawPassword"
}
Returns a JWT token that will be used to authenticate other requests
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