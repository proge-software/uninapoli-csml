package tgconf

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	//TokenKey Telegram API Key
	TokenKey = "WSS_TOKEN"
	//PollerTimeMsKey Telegram poller interval
	PollerTimeMsKey = "WSS_POLLERTIME_MS"
	//DefaultPollerTime default value for poller interval
	DefaultPollerTime = 10 * time.Second
	//MinimumPollerTime minimum value for poller interval
	MinimumPollerTime = 1 * time.Second
)

//GetToken retrieves the Telegram token from Envs
func GetToken() (*string, error) {
	if token := os.Getenv(TokenKey); token != "" {
		return &token, nil
	}
	return nil, fmt.Errorf(`required env key "%s" not set`, TokenKey)
}

//GetPollerInterval retrieves Telegram's poller duration from Envs
func GetPollerInterval() (*time.Duration, error) {
	v := os.Getenv(PollerTimeMsKey)
	if v == "" {
		duration := MinimumPollerTime
		return &duration, nil
	}

	a, err := strconv.Atoi(v)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid value for poller time (env key %s): %s. It must be a positive integer bigger than %v",
			PollerTimeMsKey, v, MinimumPollerTime)
	}

	duration := time.Duration(a)
	if duration <= MinimumPollerTime {
		return nil, fmt.Errorf(
			"invalid value for poller time (env key %s): %s. It must be a positive integer bigger than %v",
			PollerTimeMsKey, v, MinimumPollerTime)
	}

	return &duration, nil
}
