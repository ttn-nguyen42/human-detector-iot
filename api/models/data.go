package models

type SensorData struct {
	HeatLevel  int     `json:"heat_level"`
	LightLevel int     `json:"light_level"`
	DeviceId   string  `json:"device_id"`
	Timestamp  float64 `json:"timestamp"`
}
