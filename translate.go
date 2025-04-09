package deeplx_translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

type TranslationResultV2 struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func (t *Translator) TranslateText(text any, targetLang string, opts ...TranslateOption) (string, error) {
	switch t.version {
	case VersionV1:
		v, err := textToString(text)
		if err != nil {
			return "", err
		}
		resp, err := t.TranslateTextV1(v, targetLang, opts...)
		if err != nil {
			return "", err
		}
		return resp.Data, nil
	case VersionV2:
		v, err := textToStringSlice(text)
		if err != nil {
			return "", err
		}
		resp, err := t.TranslateTextV2(v, targetLang, opts...)
		if err != nil {
			return "", err
		}
		sb := &strings.Builder{}
		for _, tl := range resp.Translations {
			sb.WriteString(tl.Text)
		}
		return sb.String(), nil
	default:
		return "", fmt.Errorf("invalid API version: %d", t.version)
	}
}

func (t *Translator) TranslateTextV1(text string, targetLang string, opts ...TranslateOption) (*TranslationResultV1, error) {
	resp, err := t.translateRequest(text, targetLang, opts...)
	if err != nil {
		return nil, err
	}
	if result, ok := resp.(*TranslationResultV1); ok {
		return result, nil
	}
	return nil, fmt.Errorf("invalid response type: %T", resp)
}

func (t *Translator) TranslateTextV2(text []string, targetLang string, opts ...TranslateOption) (*TranslationResultV2, error) {
	resp, err := t.translateRequest(text, targetLang, opts...)
	if err != nil {
		return nil, err
	}
	if result, ok := resp.(*TranslationResultV2); ok {
		return result, nil
	}
	return nil, fmt.Errorf("invalid response type: %T", resp)
}

func (t *Translator) translateRequest(text any, targetLang string, opts ...TranslateOption) (any, error) {
	const (
		endpoint = "translate"
		method   = http.MethodPost
	)

	data := struct {
		Text       any    `json:"text"`
		TargetLang string `json:"target_lang"`

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
	res, err := t.callAPI(method, endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	// Parse response
	var response = map[Version]any{
		VersionV1: &TranslationResultV1{},
		VersionV2: &TranslationResultV2{},
	}[t.version]

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
