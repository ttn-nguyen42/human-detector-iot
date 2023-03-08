package models

type SensorData struct {
	Temperature int     `json:"temp"`
	Humidity    int     `json:"humidity"`
	Detected    bool    `json:"detected"`
	DeviceId    string  `json:"device_id"`
	Timestamp   float64 `json:"timestamp"`
}
