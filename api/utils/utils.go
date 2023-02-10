package utils

import (
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