package wssface

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"
	"github.com/Azure/go-autorest/autorest"
)

//FaceServiceClient client for the Azure Face Service
type FaceServiceClient struct {
	client *face.Client
}

//NewFaceServiceClient FaceServiceClient constructor
func NewFaceServiceClient(conf *Configuration) *FaceServiceClient {
	if !conf.IsValid() {
		return nil
	}

	// Client used for Detect Faces, Find Similar, and Verify examples.
	client := face.NewClient(conf.FaceEndpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(conf.FaceSubscription)

	return &FaceServiceClient{
		client: &client,
	}
}

//InvokeFace invokes the Face APIs with the provided photo
func (s *FaceServiceClient) InvokeFace(faceContext context.Context, photo io.ReadCloser) (*FaceResult, error) {
	detectSingleFaces, err := s.callFaceService(faceContext, photo)
	if err != nil {
		return nil, err
	}

	dFaces := *detectSingleFaces.Value
	faces := make([]FaceDetails, len(dFaces))
	for idx, dFace := range dFaces {
		fd := newFaceDetails(&dFace)
		faces[idx] = *fd
	}

	return &FaceResult{Faces: faces}, nil
}

func (s *FaceServiceClient) callFaceService(faceContext context.Context, photo io.ReadCloser) (face.ListDetectedFace, error) {
	// Use recognition model 2 for feature extraction. Recognition model 1 is used to simply recogize faces.
	recognitionModel02, detectionModel := face.Recognition02, face.Detection01
	// Array types chosen for the attributes of Face
	attributes := []face.AttributeType{"age", "emotion", "gender"}
	returnFaceID, returnRecognitionModel, returnFaceLandmarks := true, false, false

	return s.client.DetectWithStream(
		faceContext, photo, &returnFaceID, &returnFaceLandmarks,
		attributes, recognitionModel02, &returnRecognitionModel, detectionModel)
}

// FaceResult result of the FaceAPI
type FaceResult struct {
	Faces []FaceDetails
}

//FaceDetails Details about one face
type FaceDetails struct {
	Age       float64
	Sentiment Emotion
	Gender    string
}

func newFaceDetails(df *face.DetectedFace) *FaceDetails {
	attributes := df.FaceAttributes
	age := *attributes.Age
	gender := string(attributes.Gender)
	emotion := getEmotion(attributes)

	return &FaceDetails{
		Age:       age,
		Gender:    gender,
		Sentiment: emotion,
	}
}
