package main

import (
	"fmt"
	"log"

	"github.com/proge-software/uninapoli-csml-csbot/internal/tgconf"
	bot "github.com/proge-software/uninapoli-csml-csbot/internal/visionbot"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/envext"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssface"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssformrecognizer"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssvision"
)

// AppEnvKey Environment variable key where is stored the environment to use for the app
// execution. default is `development`
const AppEnvKey = "WSSBOT_ENV"

func main() {
	log.Println("Starting Vision telegram bot")

	envext.LoadDotenvs(AppEnvKey) // load the env vars from .env file

	conf, err := getConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	fbot, err := bot.New(*conf)
	if err != nil {
		log.Printf("can not instantiate bot: %v", err)
	}

	log.Println("Vision Telegram bot configured")

	// start bot
	fbot.Start()
}

func getConfigurations() (*bot.Configuration, error) {
	c := bot.Configuration{}

	{ // Telegram Configuration
		conf, err := tgconf.GetConfigurationsFromEnv()
		if err != nil {
			return nil, fmt.Errorf("error retrieving telegram configuration: %v", err)
		}
		c.TelegramConf = *conf
	}
	{ // Vision: Face
		faceConf, err := wssface.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving face service configuration: %v", err)
		}
		c.FaceConf = faceConf
	}
	{ // Vision API
		visionConf, err := wssvision.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving vision service configuration: %v", err)
		}
		c.VisionConf = visionConf
	}
	{ // Vision: Form Recognizer
		formRecognizerConf, err := wssformrecognizer.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving form recognizer service configuration: %v", err)
		}
		c.FormRecognizerConf = formRecognizerConf
	}

	return &c, nil
}
