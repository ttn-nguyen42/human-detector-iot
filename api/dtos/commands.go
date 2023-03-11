package dtos

const (
	ACTCheckActive string = "is-active"
)

type MQTTActivityMessage struct {
	ActionId string `json:"action_id"`
	Action string `json:"action"`
	Payload string` json:"payload"`
}

type MQTTActivityResponse struct {
	ActionId string `json:"action_id"`
	Result string `json:"result"`
	Payload string `json:"payload"`
}