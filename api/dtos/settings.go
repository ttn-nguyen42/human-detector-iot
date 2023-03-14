package dtos

type POSTCreateSettings struct {
	DataRateInSeconds int `json:"data_rate" copier:"data_rate" default:"3"`
}

type GETGetSettings struct {
	DeviceId           string   `json:"device_id"`
	DataRate           int      `json:"data_rate" copier:"data_rate"`
	NotificationEmails []string `json:"email" copier:"email"`
}

type POSTChangeEmail struct {
	Email []string `json:"email" copier:"email"`
}
