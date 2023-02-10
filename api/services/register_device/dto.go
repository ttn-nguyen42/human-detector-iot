package register_device

// The payload body for POST /register_device
type POSTRegisterDeviceDto struct {
	deviceId string `json:"device_id"`
}