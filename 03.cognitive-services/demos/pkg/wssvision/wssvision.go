package wssvision

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/computervision"
	"github.com/Azure/go-autorest/autorest"
)

//VisionServiceClient client for the Azure Vision Service
type VisionServiceClient struct {
	conf      Configuration
	visionCli *computervision.BaseClient
}

//NewVisionServiceClient VisionServiceClient constructor
func NewVisionServiceClient(conf *Configuration) *VisionServiceClient {
	if !conf.IsValid() {
		return nil
	}

	client := computervision.New(conf.ServiceEnpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(conf.VisionSubscription)

	return &VisionServiceClient{
		conf:      *conf,
		visionCli: &client,
	}
}

//InvokeVision invokes the Vision APIs with the provided photo
func (s *VisionServiceClient) InvokeVision(visionContext context.Context, photo io.ReadCloser) (*VisionResult, error) {
	if s == nil {
		return nil, fmt.Errorf("vision service client is not initialized")
	}

	features := []computervision.VisualFeatureTypes{computervision.VisualFeatureTypesDescription}
	res, err := s.visionCli.AnalyzeImageInStream(visionContext, photo, features, nil, "en", nil)
	if err != nil {
		return nil, err
	}

	if res.Description == nil || len(*res.Description.Captions) == 0 {
		return &VisionResult{Description: nil}, nil
	}

	caption := (*res.Description.Captions)[0]
	return &VisionResult{Description: caption.Text}, nil
}

// VisionResult result of the VisionAPI
type VisionResult struct {
	Description *string
}
