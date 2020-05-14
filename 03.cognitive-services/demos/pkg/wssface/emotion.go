package wssface

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"
)

// Emotion maps the Azure Face API result for emotion field
type Emotion string

const (
	anger     = Emotion("anger")
	contempt  = Emotion("contempt")
	disgust   = Emotion("disgust")
	fear      = Emotion("fear")
	happiness = Emotion("happiness")
	neutral   = Emotion("neutral")
	sadness   = Emotion("sadness")
	surprise  = Emotion("surprise")
)

// Adjective returns the adjective corresponding to the emotion
func (e Emotion) Adjective() string {
	switch e {
	case anger:
		return "angry"
	case contempt:
		return "contemptuous"
	case disgust:
		return "disgusted"
	case fear:
		return "fearful"
	case happiness:
		return "happy"
	case neutral:
		return "serious"
	case sadness:
		return "sad"
	case surprise:
		return "surprised"
	}

	return fmt.Sprintln(e)
}

func getEmotion(faceAttributes *face.Attributes) Emotion {
	emotionStruct := faceAttributes.Emotion

	var emotionMap map[string]float64
	result, _ := json.Marshal(emotionStruct)
	json.Unmarshal(result, &emotionMap)

	// Find the emotion with the highest score (confidence level). Range is 0.0 - 1.0.
	var highest float64
	var emotion string
	for name, value := range emotionMap {
		if value > highest {
			emotion, highest = name, value
		}
	}

	return Emotion(emotion)
}
