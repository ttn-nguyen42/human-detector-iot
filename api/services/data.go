package services

import (
	"fmt"
	"iot_api/models"
	"iot_api/network"
	"iot_api/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DataService interface {
	RetrieveSensorDataStream(deviceId string, quit <-chan bool) (func(), <-chan models.SensorData)
}

type dataService struct {
	Topic string
}

func NewDataService() DataService {
	return &dataService{
		Topic: "yolobit/data/sensor",
	}
}

func (s *dataService) RetrieveSensorDataStream(deviceId string, quit <-chan bool) (func(), <-chan models.SensorData) {
	res := make(chan models.SensorData)
	topic := fmt.Sprintf("%v/%v", s.Topic, deviceId)
	if !utils.IsTestMode() {
		/* If in production mode, subscribes to MQTT server */
		go func() {
			network.GetClient().Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
				// Is a duplicate
				if m.Duplicate() {
					return
				}
				var parsedPayload models.SensorData
				err := models.StructifyBytes(m.Payload(), &parsedPayload)
				if err != nil {
					return
				}
				res <- parsedPayload
			})
		}()
	} else {
		/* In test mode, data is generated */
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
					res <- utils.GetRandomSensorData(deviceId)
					time.Sleep(time.Second * 2) // Send new data every 2 seconds
				}
			}
		}()
	}
	return func() {
		network.GetClient().Unsubscribe(topic)
		close(res)
	}, res
}
