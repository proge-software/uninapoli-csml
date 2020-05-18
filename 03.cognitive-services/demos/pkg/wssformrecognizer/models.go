package wssformrecognizer

// FormRecognizerResult result of the FormRecognizerAPI
type FormRecognizerResult struct {
	Status        *string `json:"status"`
	AnalyzeResult *struct {
		DocumentResults []*struct {
			Fields *map[string]*struct {
				Text *string `json:"text"`
			} `json:"fields"`
		} `json:"documentResults"`
	} `json:"analyzeResult"`
}

//IsSucceeded checks if the result status is succeeded
func (r *FormRecognizerResult) IsSucceeded() bool {
	return r != nil && *r.Status == "succeeded"
}

//Completed checks if the job has completed
func (r *FormRecognizerResult) Completed() bool {
	return r != nil && (*r.Status == "succeeded" || *r.Status == "failed")
}

func (r *FormRecognizerResult) fieldsExists() bool {
	return r != nil && r.AnalyzeResult != nil &&
		r.AnalyzeResult.DocumentResults != nil &&
		len(r.AnalyzeResult.DocumentResults) > 0
}

func (r *FormRecognizerResult) firstField(key string) *string {
	if !r.fieldsExists() {
		return nil
	}

	for _, fields := range r.AnalyzeResult.DocumentResults {
		if fields != nil && fields.Fields != nil {
			for field, value := range *fields.Fields {
				if field == key {
					if value != nil && value.Text != nil {
						return value.Text
					}
					return nil
				}
			}
		}
	}
	return nil
}

//MerchantName ...
func (r *FormRecognizerResult) MerchantName() *string {
	return r.firstField("MerchantName")
}

//MerchantAddress ...
func (r *FormRecognizerResult) MerchantAddress() *string {
	return r.firstField("MerchantAddress")
}

//TransactionDate ...
func (r *FormRecognizerResult) TransactionDate() *string {
	return r.firstField("TransactionDate")
}

//Subtotal ...
func (r *FormRecognizerResult) Subtotal() *string {
	return r.firstField("Subtotal")
}

//Total ...
func (r *FormRecognizerResult) Total() *string {
	return r.firstField("Total")
}
