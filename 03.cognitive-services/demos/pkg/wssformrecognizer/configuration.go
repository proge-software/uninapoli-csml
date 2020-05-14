package wssformrecognizer

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

//Configuration structure for the FormRecognizer service configuration
type Configuration struct {
	FormRecognizerSubscription string
	ServiceEnpoint             string
	Retries                    struct {
		MaxAttempts int
		Interval    time.Duration
	}
}

//IsValid Checks if the configuration is valid
func (c *Configuration) IsValid() bool {
	return c != nil &&
		c.FormRecognizerSubscription != "" &&
		c.ServiceEnpoint != "" &&
		c.Retries.MaxAttempts >= 0 &&
		c.Retries.Interval >= 100*time.Millisecond
}

const (
	//FormRecognizerSubscriptionKey Azure FormRecognizer Subscription env key
	FormRecognizerSubscriptionKey = "WSS_FORMRECOGNIZER_SUBSCRIPTION_KEY"
	//FormRecognizerEndpointKey Azure FormRecognizer Endpoint env key
	FormRecognizerEndpointKey = "WSS_FORMRECOGNIZER_ENDPOINT"
	//FormRecognizerRetriesMaxAttempts Azure FormRecognizer Endpoint env key
	//for the num of max attempts to retrieve the response from the service
	FormRecognizerRetriesMaxAttempts = "WSS_FORMRECOGNIZER_RETRIES_MAXATTEMPTS"
	//FormRecognizerRetryInterval Azure FormRecognizer Endpoint env key
	//for the time to wait between each call in milliseconds
	FormRecognizerRetryInterval = "WSS_FORMRECOGNIZER_RETRIES_INTERVAL"
)

//BuildConfigurationFromEnvs builds the configuration from env variables
func BuildConfigurationFromEnvs() (*Configuration, error) {
	sub, err := getSubscriptionFromEnv()
	if err != nil {
		return nil, err
	}

	endStr, err := getEndpointFromEnv()
	if err != nil {
		return nil, err
	}

	maxAtt, err := getMaxRetriesFromEnv()
	if err != nil {
		return nil, err
	}

	retryInterval, err := getRetryIntervalFromEnv()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		FormRecognizerSubscription: *sub,
		ServiceEnpoint:             *endStr,
		Retries: struct {
			MaxAttempts int
			Interval    time.Duration
		}{
			MaxAttempts: *maxAtt,
			Interval:    *retryInterval,
		},
	}, nil
}

func getSubscriptionFromEnv() (*string, error) {
	if v := os.Getenv(FormRecognizerSubscriptionKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("formrecognizer subscription env key (%s) not set", FormRecognizerSubscriptionKey)
}

func getEndpointFromEnv() (*string, error) {
	if v := os.Getenv(FormRecognizerEndpointKey); v != "" {
		return &v, nil
	}

	return nil, fmt.Errorf("formrecognizer enpoint env key (%s) not set", FormRecognizerEndpointKey)
}

func getMaxRetriesFromEnv() (*int, error) {
	v := os.Getenv(FormRecognizerRetriesMaxAttempts)
	if v == "" {
		return nil, fmt.Errorf("formrecognizer max retries env key (%s) not set", FormRecognizerRetriesMaxAttempts)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, fmt.Errorf("can not convert formrecognizer max retries to int: %v", err)
	}

	return &i, nil
}

func getRetryIntervalFromEnv() (*time.Duration, error) {
	v := os.Getenv(FormRecognizerRetryInterval)
	if v == "" {
		return nil, fmt.Errorf("formrecognizer retry interval env key (%s) not set", FormRecognizerRetryInterval)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, fmt.Errorf("can not convert formrecognizer retry interval to int: %v", err)
	}

	d := time.Duration(i) * time.Millisecond
	return &d, nil
}
