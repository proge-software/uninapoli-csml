package wssmoderator

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/contentmoderator"
	"github.com/Azure/go-autorest/autorest"
)

//ContentModeratorServiceClient client for the Azure ContentModerator Service
type ContentModeratorServiceClient struct {
	conf              Configuration
	textModeratorCli  *contentmoderator.TextModerationClient
	imageModeratorCli *contentmoderator.ImageModerationClient
}

//NewContentModeratorServiceClient ContentModeratorServiceClient constructor
func NewContentModeratorServiceClient(conf *Configuration) *ContentModeratorServiceClient {
	if !conf.IsValid() {
		return nil
	}

	authorizer := autorest.NewCognitiveServicesAuthorizer(conf.ContentModeratorSubscription)

	textModerator := contentmoderator.NewTextModerationClient(conf.ServiceEnpoint)
	textModerator.Authorizer = authorizer

	imageModerator := contentmoderator.NewImageModerationClient(conf.ServiceEnpoint)
	imageModerator.Authorizer = authorizer

	return &ContentModeratorServiceClient{
		conf:              *conf,
		textModeratorCli:  &textModerator,
		imageModeratorCli: &imageModerator,
	}
}

//InvokeContentModeratorText invokes the ContentModerator APIs with the provided message
func (s *ContentModeratorServiceClient) InvokeContentModeratorText(ctx context.Context, message string) (*ContentModeratorTextResult, error) {
	if s == nil {
		return nil, fmt.Errorf("content moderator service client is not initialized")
	}

	msgReadCloser := ioutil.NopCloser(strings.NewReader(message))
	res, err := s.textModeratorCli.ScreenText(ctx, "text/plain", msgReadCloser, "", nil, nil, "", nil)
	if err != nil {
		return nil, err
	}

	if res.Terms == nil {
		return &ContentModeratorTextResult{}, nil
	}

	var badWords []string
	if lt := len(*res.Terms); lt > 0 {
		badWords = make([]string, lt)
		for i, term := range *res.Terms {
			badWords[i] = *term.Term
		}
	}
	return &ContentModeratorTextResult{BadWords: badWords}, nil
}

//InvokeContentModeratorPhoto invokes the ContentModerator APIs with the provided photo
func (s *ContentModeratorServiceClient) InvokeContentModeratorPhoto(ctx context.Context, photo io.ReadCloser) (*ContentModeratorPhotoResult, error) {
	if s == nil {
		return nil, fmt.Errorf("content moderator service client is not initialized")
	}

	res, err := s.imageModeratorCli.EvaluateFileInput(ctx, photo, nil)
	if err != nil {
		return nil, err
	}

	adult := res.IsImageAdultClassified != nil && *res.IsImageAdultClassified
	racy := res.IsImageRacyClassified != nil && *res.IsImageRacyClassified
	return &ContentModeratorPhotoResult{Adult: adult, Racy: racy}, nil
}

// ContentModeratorTextResult result of the ContentModeratorAPI for the evaluation of a message
type ContentModeratorTextResult struct {
	BadWords []string
}

// ContentModeratorPhotoResult result of the ContentModeratorAPI for the evaluation of a photo
type ContentModeratorPhotoResult struct {
	Adult bool
	Racy  bool
}
