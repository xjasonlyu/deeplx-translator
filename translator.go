package deeplx_translator

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	deeplProAPIURLv2  = "https://api.deepl.com/v2"
	deeplFreeAPIURLv2 = "https://api-free.deepl.com/v2"
)

type Translator struct {
	client  HTTPClient
	baseURL string
	authKey string
	version Version
}

// TranslatorOption is a functional option for configuring the Translator.
type TranslatorOption func(*Translator)

// WithBaseURL allows overriding the default base API url.
func WithBaseURL(baseURL string) TranslatorOption {
	return func(t *Translator) {
		t.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithVersion allows specifying the deepl API version.
func WithVersion(version Version) TranslatorOption {
	return func(t *Translator) {
		t.version = version
	}
}

// WithHTTPClient allows overriding the default http client.
func WithHTTPClient(c HTTPClient) TranslatorOption {
	return func(t *Translator) {
		t.client = c
	}
}

// NewTranslator creates a new translator.
func NewTranslator(authKey string, opts ...TranslatorOption) *Translator {
	// Determine default base url based on the auth key.
	var baseURL string
	if isFreeAccountAuthKey(authKey) {
		baseURL = deeplFreeAPIURLv2
	} else {
		baseURL = deeplProAPIURLv2
	}

	// Set up with default http client.
	t := &Translator{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		authKey: authKey,
	}
	t.applyOptions(opts...)

	// Auto infer API version from base url.
	if !t.version.IsValid() {
		t.version = inferVersionFromBaseURL(t.baseURL)
	}

	return t
}

// applyOptions applies the supplied functional options to the Translator.
func (t *Translator) applyOptions(opts ...TranslatorOption) {
	for _, option := range opts {
		option(t)
	}
}

// callAPI calls the supplied API endpoint with the provided parameters and returns the response.
func (t *Translator) callAPI(method string, endpoint string, headers http.Header, body io.Reader) (*http.Response, error) {
	apiURL, err := url.JoinPath(t.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error joining API url: %w", err)
	}

	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if t.authKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.authKey))
	}
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}

	return t.client.Do(req)
}

// isFreeAccountAuthKey determines whether the supplied auth key belongs to a Free account.
func isFreeAccountAuthKey(authKey string) bool {
	return strings.HasSuffix(authKey, ":fx")
}
