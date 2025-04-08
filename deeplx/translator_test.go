package deeplx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// mockHTTPClient is a simple mock of HTTPClient interface.
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TranslatorTestSuite defines the test suite for the Translator.
type TranslatorTestSuite struct {
	suite.Suite
	defaultAuthKey string
	freeAuthKey    string
}

// Run the test suite.
func TestTranslatorTestSuite(t *testing.T) {
	suite.Run(t, new(TranslatorTestSuite))
}

// SetupTest initializes shared values for tests.
func (s *TranslatorTestSuite) SetupTest() {
	s.defaultAuthKey = "pro-auth-key"
	s.freeAuthKey = "free-auth-key:fx"
}

// TestNewTranslator_DefaultClient verifies the creation of Translator with default settings.
func (s *TranslatorTestSuite) TestNewTranslator_DefaultClient() {
	tr, err := NewTranslator(s.defaultAuthKey)
	s.NoError(err)
	s.NotNil(tr)
	s.Equal(ServerURLPro, tr.serverURL)
	s.Equal(s.defaultAuthKey, tr.authKey)
}

// TestNewTranslator_FreeAccount verifies that free account auth key sets correct server URL.
func (s *TranslatorTestSuite) TestNewTranslator_FreeAccount() {
	tr, err := NewTranslator(s.freeAuthKey)
	s.NoError(err)
	s.Equal(ServerURLFree, tr.serverURL)
}

// TestWithServerURL verifies the functional option WithServerURL.
func (s *TranslatorTestSuite) TestWithServerURL() {
	opt := WithServerURL("https://custom.example.com/path")
	tr := &Translator{}
	err := opt(tr)
	s.NoError(err)
	s.Equal("https://custom.example.com", tr.serverURL)
}

// TestWithServerURL_InvalidURL checks error behavior on malformed URL.
func (s *TranslatorTestSuite) TestWithServerURL_InvalidURL() {
	opt := WithServerURL("::bad-url")
	tr := &Translator{}
	err := opt(tr)
	s.Error(err)
}

// TestWithHTTPClient verifies the functional option WithHTTPClient.
func (s *TranslatorTestSuite) TestWithHTTPClient() {
	mockClient := &mockHTTPClient{}
	opt := WithHTTPClient(mockClient)
	tr := &Translator{}
	err := opt(tr)
	s.NoError(err)
	s.Equal(mockClient, tr.client)
}

// TestCallAPI_Success simulates a successful API call.
func (s *TranslatorTestSuite) TestCallAPI_Success() {
	mockResp := "translated"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal(http.MethodPost, r.Method)
		s.Equal("DeepL-Auth-Key "+s.defaultAuthKey, r.Header.Get("Authorization"))
		s.Equal("/v2/translate", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResp))
	}))
	defer server.Close()

	tr, err := NewTranslator(s.defaultAuthKey,
		WithHTTPClient(&http.Client{}),
		WithServerURL(server.URL),
	)
	s.NoError(err)

	res, err := tr.callAPI(http.MethodPost, "/v2/translate", http.Header{
		"Content-Type": []string{"application/json"},
	}, strings.NewReader(`{"text": "Hello"}`))

	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	s.Equal(mockResp, string(body))
}

// TestCallAPI_RetriableError simulates a path join error.
func (s *TranslatorTestSuite) TestCallAPI_InvalidJoinPath() {
	tr, err := NewTranslator(s.defaultAuthKey,
		WithHTTPClient(&http.Client{}),
	)
	s.NoError(err)

	// Inject invalid serverURL to trigger url.JoinPath error
	trWithInvalidURL := *tr
	trWithInvalidURL.serverURL = "://invalid-url"

	_, err = trWithInvalidURL.callAPI(http.MethodPost, "/endpoint", nil, nil)
	s.Error(err)
	s.ErrorContains(err, "error joining API url")
}

// TestCallAPI_RetriableError simulates a retryable error.
func (s *TranslatorTestSuite) TestCallAPI_RetriableError() {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr, err := NewTranslator(s.defaultAuthKey,
		WithHTTPClient(&http.Client{}),
		WithServerURL(server.URL),
	)
	s.NoError(err)

	res, err := tr.callAPI("GET", "/v2/check", nil, nil)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.GreaterOrEqual(attempts, 3)
}

// TestCallAPI_NonRetriableError returns immediate error for bad request.
func (s *TranslatorTestSuite) TestCallAPI_NonRetriableError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	tr, err := NewTranslator(s.defaultAuthKey,
		WithHTTPClient(&http.Client{}),
		WithServerURL(server.URL),
	)
	s.NoError(err)

	res, err := tr.callAPI("GET", "/v2/bad", nil, nil)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, res.StatusCode)
}

// TestIsFreeAccountAuthKey verifies detection of free auth keys.
func (s *TranslatorTestSuite) TestIsFreeAccountAuthKey() {
	s.True(isFreeAccountAuthKey("abc:fx"))
	s.False(isFreeAccountAuthKey("abc123"))
}
