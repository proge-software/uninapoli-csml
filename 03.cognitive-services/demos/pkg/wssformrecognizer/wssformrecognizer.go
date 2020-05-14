package wssformrecognizer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	suffix                  = "formrecognizer/v2.0-preview/prebuilt/receipt/analyze"
	operationLocationHeader = "Operation-Location"
)

//FormRecognizerServiceClient client for the Azure FormRecognizer Service
type FormRecognizerServiceClient struct {
	conf              Configuration
	formRecognizerCli *http.Client
}

//NewFormRecognizerServiceClient FormRecognizerServiceClient constructor
func NewFormRecognizerServiceClient(conf *Configuration) *FormRecognizerServiceClient {
	if !conf.IsValid() {
		return nil
	}

	client := http.Client{Timeout: 10 * time.Second}

	return &FormRecognizerServiceClient{
		conf:              *conf,
		formRecognizerCli: &client,
	}
}

//InvokeFormRecognizer invokes the FormRecognizer APIs with the provided photo
func (s *FormRecognizerServiceClient) InvokeFormRecognizer(ctx context.Context, photo io.ReadCloser) (*FormRecognizerResult, error) {
	if s == nil {
		return nil, fmt.Errorf("content form recognizer service client is not initialized")
	}

	req, err := s.newProcessPhotoRequest(ctx, photo)
	if err != nil {
		return nil, err
	}

	resp, err := s.doFormRecognizerRequest(req)
	if err != nil {
		return nil, err
	}

	log.Printf("Form analysis succeded: %v", resp.Header)

	req, err = s.newGetResultRequest(ctx, resp)
	if err != nil {
		return nil, err
	}

	for i := 0; i < s.conf.Retries.MaxAttempts; i++ {
		res, err := s.doGetResultRequest(req)
		if err != nil {
			return nil, err
		}
		if res.Completed() {
			return res, err
		}
		time.Sleep(s.conf.Retries.Interval)
	}

	return nil, fmt.Errorf(
		"no successfull response obtained after %v tentative (interval: %v)",
		s.conf.Retries.MaxAttempts, s.conf.Retries.Interval)
}

func (s *FormRecognizerServiceClient) newProcessPhotoRequest(ctx context.Context, photo io.ReadCloser) (*http.Request, error) {
	url, err := url.Parse(s.conf.ServiceEnpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating post photo request for form recognizer (url builder): %v", err)
	}
	url.Path = suffix

	postURL := url.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, photo)
	if err != nil {
		return nil, fmt.Errorf("error creating post photo request for form recognizer: %v", err)
	}

	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Ocp-Apim-Subscription-Key", s.conf.FormRecognizerSubscription)

	return req, nil
}

func (s *FormRecognizerServiceClient) doFormRecognizerRequest(request *http.Request) (*http.Response, error) {
	response, err := s.formRecognizerCli.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("form analysis failed; can not read body: %v", err)
		}
		return nil, fmt.Errorf("form analysis failed: %v", body)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}
	log.Printf("response body: %s", body)

	return response, nil
}

func (s *FormRecognizerServiceClient) getResultURL(resp *http.Response) (*string, error) {
	opLocs := resp.Header[operationLocationHeader]
	if len(opLocs) == 0 || opLocs[0] == "" {
		return nil, fmt.Errorf("form recognizer response do not have header %s: %v", operationLocationHeader, resp.Header)
	}
	return &opLocs[0], nil
}

func (s *FormRecognizerServiceClient) newGetResultRequest(ctx context.Context, resp *http.Response) (*http.Request, error) {
	getURL, err := s.getResultURL(resp)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating get result request for form recognizer: %v", err)
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", s.conf.FormRecognizerSubscription)

	return req, nil
}

func (s *FormRecognizerServiceClient) doGetResultRequest(request *http.Request) (*FormRecognizerResult, error) {
	response, err := s.formRecognizerCli.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("form analysis failed; can not read body: %v", err)
	}
	log.Printf("read body: %s", string(body))

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("form analysis failed: %v", body)
	}

	var result FormRecognizerResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling json response: %v", err)
	}

	return &result, nil
}
