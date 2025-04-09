package deeplx_translator

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateText(t *testing.T) {
	var (
		deeplAPIKey  = os.Getenv("DEEPL_API_KEY")
		deeplxAPIKey = os.Getenv("DEEPLX_API_KEY")
		deeplxAPIURL = os.Getenv("DEEPLX_API_URL")
	)

	for _, test := range []struct {
		name           string
		apiKey, apiURL string
		version        Version
	}{
		{"DeepL Official API", deeplAPIKey, "", 0},
		{"DeepLX Free API", deeplxAPIKey, deeplxAPIURL + "/", 0},
		// {"DeepLX Pro API", deeplxAPIKey, deeplxAPIURL + "/v1", 0},
		{"DeepLX Official API", deeplxAPIKey, deeplxAPIURL + "/v2", VersionV2},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.apiKey == "" {
				t.SkipNow()
			}

			opts := []TranslatorOption{
				WithHTTPClient(http.DefaultClient),
			}
			if test.apiURL != "" {
				opts = append(opts, WithBaseURL(test.apiURL))
			}
			if test.version.IsValid() {
				opts = append(opts, WithVersion(test.version))
			}
			translator := NewTranslator(test.apiKey, opts...)

			for _, unit := range []struct {
				text     any
				from, to string
			}{
				{`Oh yeah! I'm a translator!`, "", "zh"},
				{`Oh yeah! I'm a translator!`, "", "zh-Hant"},
				{`Oh yeah! I'm a translator!`, "", "ja"},
				{[]string{`Oh yeah! I'm a translator!`}, "", "de"},
				{[]string{`Oh yeah! I'm a translator!`}, "en", "fr"},
			} {
				result, err := translator.TranslateText(
					unit.text, unit.to,
					WithSourceLang(unit.from),
				)
				if assert.NoError(t, err) {
					t.Log(result)
				}
			}
		})
	}
}
