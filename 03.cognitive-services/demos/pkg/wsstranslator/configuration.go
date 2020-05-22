package wsstranslator

import (
	"fmt"
	"os"
)

//Configuration structure for the Translator service configuration
type Configuration struct {
	TranslatorRegion       string
	TranslatorSubscription string
	ServiceEnpoint         string
}

const (
	//TranslatorSubscriptionKey Azure Translator Subscription env key
	TranslatorSubscriptionKey = "WSS_TRANSLATOR_SUBSCRIPTION_KEY"
	//TranslatorEndpointKey Azure Translator Endpoint env key
	TranslatorEndpointKey = "WSS_TRANSLATOR_ENDPOINT"
	//TranslatorRegionKey Azure Translator Region env key
	TranslatorRegionKey = "WSS_TRANSLATOR_REGION"
)

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil &&
		c.TranslatorSubscription != "" &&
		c.ServiceEnpoint != "" &&
		c.TranslatorRegion != ""
}

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	subscription, err := getTranslatorSubscription()
	if err != nil {
		return nil, err
	}

	endpoint, err := getTranslatorEndpoint()
	if err != nil {
		return nil, err
	}

	region, err := getTranslatorRegion()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		TranslatorRegion:       *region,
		TranslatorSubscription: *subscription,
		ServiceEnpoint:         *endpoint,
	}, nil
}

func getTranslatorSubscription() (*string, error) {
	if v := os.Getenv(TranslatorSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("translator subscription env key (%s) not set", TranslatorSubscriptionKey)
}

func getTranslatorEndpoint() (*string, error) {
	if v := os.Getenv(TranslatorEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("translator subscription env key (%s) not set", TranslatorEndpointKey)
}

func getTranslatorRegion() (*string, error) {
	if v := os.Getenv(TranslatorRegionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("translator region env key (%s) not set", TranslatorRegionKey)
}
