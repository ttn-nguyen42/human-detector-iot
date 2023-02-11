package dtos

// The payload body for POST /register_device
type POSTRegisterDeviceDto struct {
	DeviceId string `json:"device_id"`
}

type POSTRegisterDeviceResponse struct {
	Password string `json:"password"`
}