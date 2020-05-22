package tgconf

import (
	"fmt"
	"time"
)

//Configuration Bot's Configuration
type Configuration struct {
	PollerInterval time.Duration
	Token          string
}

//GetConfigurationsFromEnv ...
func GetConfigurationsFromEnv() (*Configuration, error) {
	c := Configuration{}

	{ // Telegram Configuration
		pollerInterval, err := GetPollerInterval()
		if err != nil {
			return nil, fmt.Errorf("error retrieving poller interval: %v", err)
		}
		c.PollerInterval = *pollerInterval

		token, err := GetToken()
		if err != nil {
			return nil, fmt.Errorf("error retrieving Telegram token: %v", err)
		}
		c.Token = *token
	}
	return &c, nil
}
