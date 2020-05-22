package wsstranslator

import (
	"context"
)

//InvokeLanguages invokes the Languages APIs to obtain the list of supported languages
func (s *TranslatorServiceClient) InvokeLanguages(LanguagesContext context.Context) (*LanguagesResult, error) {
	res, err := s.translatorCli.Languages(LanguagesContext, []string{`translation`}, `en`, "")
	if err != nil {
		return nil, err
	}

	if res.Translation == nil || *res.Translation == nil {
		return &LanguagesResult{}, nil
	}

	translations := *res.Translation
	languages, counter := make([]Language, len(translations)), 0
	for k, v := range translations {
		var name string
		if v.NativeName != nil {
			name = *v.NativeName
		} else if v.Name != nil {
			name = *v.Name
		} else {
			name = ""
		}

		languages[counter] = Language{Code: k, Name: name}
		counter++
	}
	return &LanguagesResult{Languages: languages}, nil
}

//LanguagesResult result of the LanguagesAPI
type LanguagesResult struct {
	Languages []Language
}

//Language describes a single language
type Language struct {
	Code string
	Name string
}
