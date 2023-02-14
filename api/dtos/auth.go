package dtos

type POSTLoginRequest struct {
	DeviceId string `json:"device_id"`
	Password string `json:"password"`
}

type POSTLoginResponse struct {
	Token string `json:"token"`
}