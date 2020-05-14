package wssvision

import (
	"fmt"
	"os"
)

//Configuration structure for the Vision service configuration
type Configuration struct {
	VisionSubscription string
	ServiceEnpoint     string
}

const (
	//VisionSubscriptionKey Azure Vision Subscription env key
	VisionSubscriptionKey = "WSS_VISION_SUBSCRIPTION_KEY"
	//VisionEndpointKey Azure Vision Endpoint env key
	VisionEndpointKey = "WSS_VISION_ENDPOINT"
)

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil && c.VisionSubscription != "" && c.ServiceEnpoint != ""
}

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getVisionSubscription()
	if err != nil {
		return nil, err
	}

	endStr, err := getVisionEndpoint()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		VisionSubscription: *sub,
		ServiceEnpoint:     *endStr,
	}, nil
}

func getVisionSubscription() (*string, error) {
	if v := os.Getenv(VisionSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("vision subscription env key (%s) not set", VisionSubscriptionKey)
}

func getVisionEndpoint() (*string, error) {
	if v := os.Getenv(VisionEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("vision subscription env key (%s) not set", VisionEndpointKey)
}
