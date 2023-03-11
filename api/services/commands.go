package services

import (
	"fmt"
	"iot_api/custom"
	"iot_api/dtos"
	"iot_api/models"
	"iot_api/network"
	"iot_api/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type CommandService interface {
	sendCommand(deviceId string, kind string, payload string) error
	SendStatusCheck(deviceId string) error
	SendDataRateUpdate(deviceId string, settings dtos.POSTDataRateRequest) error
}

type commandService struct {
	ActivityTopic string
	ResponseTopic string
}

func NewCommandService(activityTopic string, responseTopic string) CommandService {
	return &commandService{
		ActivityTopic: activityTopic,
		ResponseTopic: responseTopic,
	}
}

/*
 * Sends a status check to gateway
 */
func (s *commandService) SendStatusCheck(deviceId string) error {
	return s.sendCommand(deviceId, dtos.ACTCheckActive, "")
}

/*
 * Sends a data rate settings to gateway
 */
func (s *commandService) SendDataRateUpdate(deviceId string, settings dtos.POSTDataRateRequest) error {
	str, err := models.Stringify(settings)
	if err != nil {
		return custom.NewInvalidFormatError("Provided settings in invalid")
	}
	return s.sendCommand(deviceId, dtos.ACTChangeDataRate, str)
}

func (s *commandService) sendCommand(deviceId string, kind string, payload string) error {
	res := make(chan dtos.MQTTActivityResponse, 1)
	quit := make(chan bool, 1)
	actionId := utils.GetRandomUUID()
	reqTopic := fmt.Sprintf("%v/%v", s.ActivityTopic, deviceId)
	resTopic := fmt.Sprintf("%v/%v", s.ResponseTopic, deviceId)
	req := dtos.MQTTActivityMessage{
		ActionId: actionId,
		Action:   kind,
		Payload:  payload,
	}
	/* Sends request */
	if !utils.IsTestMode() {
		go func() {
			select {
			case <-quit:
				return
			default:
				reqStr, err := models.Stringify(req)
				if err != nil {
					logrus.Error(err)
					quit <- true
					return
				}
				token := network.GetClient().Publish(reqTopic, 1, false, reqStr)
				logrus.Infof("Sent request: %v", reqStr)
				token.Wait()
				if token.Error() != nil {
					logrus.Error(token.Error())
					quit <- true
					return
				}
				quit <- false
				return
			}
		}()
	}
	needQuit := <-quit
	if needQuit {
		return custom.NewUnableToSendMessage("")
	}
	if !utils.IsTestMode() {
		/* If in production mode, subscribes to MQTT server */
		go func() {
			network.GetClient().Subscribe(resTopic, 1, func(c mqtt.Client, m mqtt.Message) {
				// Is a duplicate
				if m.Duplicate() {
					return
				}
				var parsedPayload dtos.MQTTActivityResponse
				err := models.StructifyBytes(m.Payload(), &parsedPayload)
				if parsedPayload.ActionId == actionId {
					res <- parsedPayload
				}
				if err != nil {
					return
				}
			})
		}()
	} else {
		/* In test mode, data is generated */
		go func() {
			res <- dtos.MQTTActivityResponse{
				ActionId: actionId,
				Result:   "success",
				Payload:  "",
			}
		}()
	}
	select {
	case status := <-res:
		network.GetClient().Unsubscribe(resTopic)
		close(quit)
		close(res)
		if status.Result != "success" {
			return custom.NewInactiveGatewayError("")
		}
		logrus.Info(fmt.Sprintf("Status check: %v", status))
		return nil
	case <-time.After(6 * time.Second):
		network.GetClient().Unsubscribe(resTopic)
		close(quit)
		close(res)
		logrus.Info("Status check timed out")
		return custom.NewTimeoutError("")
	}
}
