package auths

import (
	"iot_api/custom"
	"iot_api/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var signedString string

func init() {
	signedKey, err := utils.GetJwtSignKey()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	signedString = signedKey
}

// Factory
func GenerateJwt(id string) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodEdDSA)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(8000 * time.Hour)
	claims["authorized"] = true
	claims["device"] = id

	signedToken, err := jwtToken.SignedString(signedString)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(token string) (string, error) {
	// Decorator
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodECDSA)
		if ok {
			return []byte(signedString), nil
		}
		return nil, custom.NewInvalidFormatError("Invalid JWT")
	})
	if !parsedToken.Valid {
		return "", custom.NewUnauthorizedError("")
	}
	if err != nil {
		return "", custom.NewInvalidFormatError("Invalid JWT")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", custom.NewFieldMissingError("Missing deviceId")
	}
	return claims["device"].(string), nil
}
