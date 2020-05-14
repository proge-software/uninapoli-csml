package wsssentiment

import (
	"fmt"
	"os"
)

//Configuration structure for the TextAnalytics service configuration
type Configuration struct {
	TextAnalyticsSubscription string
	ServiceEnpoint            string
}

const (
	//TextAnalyticsSubscriptionKey Azure TextAnalytics Subscription env key
	TextAnalyticsSubscriptionKey = "WSS_TEXTANALYTICS_SUBSCRIPTION_KEY"
	//TextAnalyticsEndpointKey Azure TextAnalytics Endpoint env key
	TextAnalyticsEndpointKey = "WSS_TEXTANALYTICS_ENDPOINT"
)

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil && c.TextAnalyticsSubscription != "" && c.ServiceEnpoint != ""
}

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getTextAnalyticsSubscription()
	if err != nil {
		return nil, err
	}

	endStr, err := getTextAnalyticsEndpoint()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		TextAnalyticsSubscription: *sub,
		ServiceEnpoint:            *endStr,
	}, nil
}

func getTextAnalyticsSubscription() (*string, error) {
	if v := os.Getenv(TextAnalyticsSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("textanalytics subscription env key (%s) not set", TextAnalyticsSubscriptionKey)
}

func getTextAnalyticsEndpoint() (*string, error) {
	if v := os.Getenv(TextAnalyticsEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("textanalytics subscription env key (%s) not set", TextAnalyticsEndpointKey)
}
