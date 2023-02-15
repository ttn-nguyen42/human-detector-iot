package auths

import (
	"iot_api/custom"
	"iot_api/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var signedString string
var issuerString string

func init() {
	signedKey, err := utils.GetJwtSignKey()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	signedString = signedKey
	issuerString = "Backend"
}

type customClaims struct {
	jwt.StandardClaims
	DeviceId string `json:"device_id"`
}

// Factory
func GenerateJwt(id string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &customClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    issuerString,
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	logrus.Debug(signedString)
	signedToken, err := jwtToken.SignedString([]byte(signedString))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(token string) (string, error) {
	// Decorator
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if ok {
			return []byte(signedString), nil
		}
		return nil, custom.NewInvalidFormatError("Cannot parse JWT")
	})
	if !parsedToken.Valid {
		return "", custom.NewUnauthorizedError(err.Error())
	}
	if err != nil {
		return "", custom.NewInvalidFormatError(err.Error())
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", custom.NewFieldMissingError("Missing deviceId")
	}
	return claims["device_id"].(string), nil
}
