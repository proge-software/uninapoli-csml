package wssface

import (
	"fmt"
	"os"
)

//Configuration structure for the Face service configuration
type Configuration struct {
	FaceSubscription string
	FaceEndpoint     string
}

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil && c.FaceEndpoint != "" && c.FaceSubscription != ""
}

const (
	//FaceSubscriptionKey Azure Face Subscription env key
	FaceSubscriptionKey = "WSS_FACE_SUBSCRIPTION_KEY"
	//FaceEndpointKey Azure Face Endpoint env key
	FaceEndpointKey = "WSS_FACE_ENDPOINT"
)

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getFaceSubscription()
	if err != nil {
		return nil, err
	}

	end, err := getFaceEndpoint()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		FaceSubscription: *sub,
		FaceEndpoint:     *end,
	}, nil
}

func getFaceSubscription() (*string, error) {
	if v := os.Getenv(FaceSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("face subscription env key (%s) not set", FaceSubscriptionKey)
}

func getFaceEndpoint() (*string, error) {
	if v := os.Getenv(FaceEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("face subscription env key (%s) not set", FaceEndpointKey)
}
