package deeplx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TranslationResultV1 holds the results of a text translation request of v1.
type TranslationResultV1 struct {
	Code         int      `json:"code"`
	ID           int64    `json:"id"`
	Message      string   `json:"message,omitempty"`
	Data         string   `json:"data"`
	Alternatives []string `json:"alternatives"`
	SourceLang   string   `json:"source_lang"`
	TargetLang   string   `json:"target_lang"`
	Method       string   `json:"method"`
}

// TranslationResultV2 holds the results of a text translation request of v2.
type TranslationResultV2 struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

// TranslateTextV2 translates the given text(s) into the specified target language.
func (t *Translator) translateText(text []string, targetLang string, version APIVersion, opts ...TranslateOption) ([]byte, error) {
	var endpoint string
	switch version {
	case FreeAPI:
		endpoint = "translate"
	case ProAPI:
		endpoint = "v1/translate"
	case OfficialAPI:
		endpoint = "v2/translate"
	default:
		return nil, fmt.Errorf("invalid API version: %d", version)
	}

	data := struct {
		Text       []string `json:"text"`
		TargetLang string   `json:"target_lang"`

		TranslateOptions
	}{
		Text:       text,
		TargetLang: targetLang,
	}
	if err := data.TranslateOptions.Gather(opts...); err != nil {
		return nil, fmt.Errorf("error setting translate option: %w", err)
	}

	// Setup request
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding request data: %w", err)
	}

	// Send request
	res, err := t.callAPI(http.MethodPost, endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	return io.ReadAll(res.Body)
}

func (t *Translator) TranslateText(text []string, targetLang string, opts ...TranslateOption) (string, error) {
	data, err := t.translateText(text, targetLang, OfficialAPI, opts...)
	if err != nil {
		return "", err
	}

	// Parse response
	var response TranslationResultV1
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&response); err != nil {
		return "", err
	}

	return "", nil
}

// TranslateTextV1 translates the given text(s) into the specified target language.
func (t *Translator) TranslateTextV1(text []string, targetLang string, opts ...TranslateOption) (*TranslationResultV1, error) {
	data, err := t.translateText(text, targetLang, OfficialAPI, opts...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response TranslationResultV1
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// TranslateTextV2 translates the given text(s) into the specified target language.
func (t *Translator) TranslateTextV2(text []string, targetLang string, opts ...TranslateOption) (*TranslationResultV2, error) {
	data, err := t.translateText(text, targetLang, OfficialAPI, opts...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response TranslationResultV2
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
