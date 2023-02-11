package utils

import (
	"crypto/md5"
	"encoding/hex"
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
	return "", nil
}

/*
Hashes a string using MD5
Not secure, but usable
*/
func MD5Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}