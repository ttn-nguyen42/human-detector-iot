package network

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"iot_api/utils"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var once sync.Once

// Singleton
var mqttClient mqtt.Client

func GetClient() mqtt.Client {
	if mqttClient != nil {
		return mqttClient
	}
	once.Do(func() {
		opts := getAwsMqttSettings()
		mqttClient = connect(opts)
	})
	return mqttClient
}

func connect(options *mqtt.ClientOptions) mqtt.Client {
	client := mqtt.NewClient(options)
	if utils.IsTestMode() {
		/* In test mode, connection to MQTT server is not allowed */
		logrus.Info("Currently in test mode, MQTT server is not needed")
		return client
	}
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		logrus.Fatal("AWS IoT Core MQTT connection failed: %v", token.Error().Error())
		return nil
	}
	return client
}

func getAwsMqttSettings() *mqtt.ClientOptions {
	if utils.IsTestMode() {
		logrus.Info("In test mode, skip checking AWS IoT Core credentials")
		return mqtt.NewClientOptions()
	}
	certs, err := utils.GetAWSIoTCertPaths()
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	endpoint, err := utils.GetAWSIoTEndpoint()
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	tlsSettings, err := getTlsSettings(certs)
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	clientOptions := mqtt.NewClientOptions()
	clientOptions.AddBroker(fmt.Sprintf("tls://%v:%v", endpoint, utils.AWSIoTMQTTPort))
	clientOptions.SetTLSConfig(tlsSettings)
	clientId := utils.GetRandomUUID()
	clientOptions.SetClientID(clientId)
	clientOptions.OnConnect = onConnect
	clientOptions.OnConnectionLost = onDisconnect
	return clientOptions
}

func getTlsSettings(certPaths *utils.CertsPaths) (*tls.Config, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(certPaths.RootCA)
	if err != nil {
		return nil, fmt.Errorf("invalid root CA path")
	}
	certPool.AppendCertsFromPEM(ca)
	clientCertKeyPair, err := tls.LoadX509KeyPair(certPaths.Cert, certPaths.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid cert and private key")
	}
	return &tls.Config{
		RootCAs: certPool,
		ClientAuth: tls.NoClientCert,
		ClientCAs: nil,
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{clientCertKeyPair},
	}, nil
}

var onConnect mqtt.OnConnectHandler = func(c mqtt.Client) {
	logrus.Info("AWS IoT Core MQTT connected")
}

var onDisconnect mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
	logrus.Info("AWS IoT Core MQTT connection lost: %v", err.Error())
}