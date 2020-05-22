package moderatebot

import (
	"github.com/proge-software/uninapoli-csml-csbot/internal/tgconf"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssface"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssformrecognizer"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssmoderator"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsssentiment"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsstranslator"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssvision"
)

//Configuration Bot's Configuration
type Configuration struct {
	TelegramConf       tgconf.Configuration
	FaceConf           *wssface.Configuration
	VisionConf         *wssvision.Configuration
	TextAnalyticsConf  *wsssentiment.Configuration
	ModeratorConf      *wssmoderator.Configuration
	TranslatorConf     *wsstranslator.Configuration
	FormRecognizerConf *wssformrecognizer.Configuration
}
