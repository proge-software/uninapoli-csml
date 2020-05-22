package visionbot

import (
	"github.com/proge-software/uninapoli-csml-csbot/internal/tgconf"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssface"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssformrecognizer"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssvision"
)

//Configuration Bot's Configuration
type Configuration struct {
	TelegramConf       tgconf.Configuration
	FaceConf           *wssface.Configuration
	VisionConf         *wssvision.Configuration
	FormRecognizerConf *wssformrecognizer.Configuration
}
