package wssmoderator

import (
	"fmt"
	"os"
)

//Configuration structure for the ContentModerator service configuration
type Configuration struct {
	ContentModeratorSubscription string
	ServiceEnpoint               string
}

const (
	//ContentModeratorSubscriptionKey Azure ContentModerator Subscription env key
	ContentModeratorSubscriptionKey = "WSS_CONTENTMODERATOR_SUBSCRIPTION_KEY"
	//ContentModeratorEndpointKey Azure ContentModerator Endpoint env key
	ContentModeratorEndpointKey = "WSS_CONTENTMODERATOR_ENDPOINT"
)

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil && c.ContentModeratorSubscription != "" && c.ServiceEnpoint != ""
}

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getContentModeratorSubscription()
	if err != nil {
		return nil, err
	}

	endStr, err := getContentModeratorEndpoint()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ContentModeratorSubscription: *sub,
		ServiceEnpoint:               *endStr,
	}, nil
}

func getContentModeratorSubscription() (*string, error) {
	if v := os.Getenv(ContentModeratorSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("contentmoderator subscription env key (%s) not set", ContentModeratorSubscriptionKey)
}

func getContentModeratorEndpoint() (*string, error) {
	if v := os.Getenv(ContentModeratorEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("contentmoderator subscription env key (%s) not set", ContentModeratorEndpointKey)
}
