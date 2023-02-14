package auths

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.Request.Header["Authorization"]
		if len(bearer) > 1 {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: "Authorization header must only have one value",
			})
			ctx.Abort()
			return
		}
		token := bearer[0]
		if !strings.HasPrefix(token, "Bearer ") {
			ctx.JSON(http.StatusBadRequest, MessageResponse{
				Message: "Authorization token must starts with 'Bearer '",
			})
			ctx.Abort()
			return
		}
		_, trimmed, _ := strings.Cut(token, "Bearer ")
		logrus.Debug(trimmed)
		deviceId, err := VerifyToken(trimmed)
		if err != nil {
			logrus.Debug(err.Error())
			ctx.JSON(http.StatusUnauthorized, MessageResponse{
				Message: "Invalid authorization token",
			})
			ctx.Abort()
			return
		}
		// Pass the extracted device_id to handlers
		// Do not verify this device_id again because we know that this JWT came from us
		// If we did everything right when establishing this JWT
		// no validation is required here
		ctx.Set("device_id", deviceId)
		ctx.Next()
	}
}
