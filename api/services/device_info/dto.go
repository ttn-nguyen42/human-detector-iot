package deviceInfo

// The payload body for POST /register_device
type POSTRegisterDeviceDto struct {
	DeviceId string `json:"device_id"`
}