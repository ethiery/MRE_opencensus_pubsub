package common

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	GCPProjectID string
	Topic        string
}

const gcpProjectIDKey = "GCP_PROJECT"
const topicKey = "TOPIC"

func GetConfig() (*Config, error) {
	var err error
	config := &Config{}
	if config.GCPProjectID, err = getEnv(gcpProjectIDKey); err != nil {
		return nil, err
	}
	if config.Topic, err = getEnv(topicKey); err != nil {
		return nil, err
	}
	return config, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", errors.New(fmt.Sprintf("Environment variable %s is not set or empty.", key))
	}
	return v, nil
}
