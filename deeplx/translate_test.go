package deeplx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TranslateTextTestSuite defines the suite for translation methods.
type TranslateTextTestSuite struct {
	suite.Suite
	authKey string
}

func TestTranslateTextTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTextTestSuite))
}

func (s *TranslateTextTestSuite) SetupTest() {
	s.authKey = "test-key"
}

// TestTranslateTextV2 verifies a successful translation
func (s *TranslateTextTestSuite) TestTranslateTextV2() {
	mockResponse := `{
		"translations": [
			{ "detected_source_language": "EN", "text": "Hallo" },
			{ "detected_source_language": "EN", "text": "Welt" }
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/v2/translate", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		s.Contains(string(body), `"target_lang":"DE"`)
		s.Contains(string(body), `"text":["Hello","World"]`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	tr, err := NewTranslator(s.authKey,
		WithServerURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	s.NoError(err)

	result, err := tr.TranslateTextV2([]string{"Hello", "World"}, "DE")
	s.NoError(err)
	s.Len(result.Translations, 2)
	s.Equal("Hallo", result.Translations[0].Text)
	s.Equal("EN", result.Translations[0].DetectedSourceLanguage)
}
