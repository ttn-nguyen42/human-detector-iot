package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

func init() {
	/* Get unique random character by changing seed everytime application launches */
	rand.Seed(time.Now().UnixMicro())
}

// Generates random UUID for passwords
func GetRandomUUID() string {
	id := uuid.New()
	return id.String()
}

// Get port from environment variables
func GetPort() string {
	port := os.Getenv(EnvPort)
	if len(port) == 0 {
		return "8080"
	}
	return port
}

// Get the database driver from environment variables
func GetMongoDriver() (string, error) {
	clusterId := os.Getenv(EnvMongoClusterID)
	if len(clusterId) == 0 {
		return "", fmt.Errorf("missing MongoDB cluster ID - %v", EnvMongoClusterID)
	}
	username := os.Getenv(EnvMongoUsername)
	if len(username) == 0 {
		return "", fmt.Errorf("missing MongoDB username - %v", EnvMongoUsername)
	}
	password := os.Getenv(EnvMongoPassword)
	if len(password) == 0 {
		return "", fmt.Errorf("missing MongoDB password - %v", EnvMongoPassword)
	}
	settings := "?retryWrites=true&w=majority"
	driver := fmt.Sprintf("mongodb+srv://%v:%v@%v.mongodb.net/%v", username, password, clusterId, settings)
	return driver, nil
}

/*
Hashes a string using MD5
Not secure, but usable
*/
func MD5Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

type CertsPaths struct {
	RootCA string
	PrivateKey string
	Cert string
}

func GetAWSIoTCertPaths() (*CertsPaths, error) {
	rootCaPath := os.Getenv(EnvAWSIoTRootCA)
	if len(rootCaPath) == 0 {
		return nil, fmt.Errorf("missing AWS Root CA path")
	}
	privateKeyPath := os.Getenv(EnvAWSIoTPrivateKey)
	if len(privateKeyPath) == 0 {
		return nil, fmt.Errorf("missing AWS IoT private key path")
	}
	certPath := os.Getenv(EnvAWSIoTCert)
	if len(certPath) == 0 {
		return nil, fmt.Errorf("missing AWS IoT cert path")
	}
	return &CertsPaths{
		RootCA: rootCaPath,
		PrivateKey: privateKeyPath,
		Cert: certPath,
	}, nil
}

func GetAWSIoTEndpoint() (string, error) {
	endpoint := os.Getenv(EnvAWSIoTEndpoint)
	if len(endpoint) == 0 {
		return "", fmt.Errorf("missing AWS IoT endpoint")
	}
	return endpoint, nil
}

func GetJwtSignKey() (string, error) {
	key := os.Getenv(EnvJwtSignKey)
	if len(key) == 0 {
		return "", fmt.Errorf("missing JWT sign key")
	}
	return key, nil
}

func GetPasswordHash(raw string) string {
	return MD5Hash(raw)
}