package routes

import "github.com/gin-gonic/gin"

/*
Generate a public RSA key for the IoT device
This public key then is going to be used to authenticate that device
though POST /api/backend/register_device
*/
func GETGetPublicKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}