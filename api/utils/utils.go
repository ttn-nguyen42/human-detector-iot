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
