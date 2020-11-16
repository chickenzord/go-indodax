package indodax

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client Indodax client
type Client struct {
	BaseURL       string
	PublicBaseURL string
	APIKey        string
	SecretKey     string
	ReceiveWindow time.Duration
}

// NewDefaultClient creates new default client
func NewDefaultClient(apiKey, secretKey string) *Client {
	return &Client{
		BaseURL:       "https://indodax.com/tapi",
		PublicBaseURL: "https://indodax.com/api",
		APIKey:        apiKey,
		SecretKey:     secretKey,
		ReceiveWindow: 5 * time.Second,
	}
}

// CallPublic generic function to call public Indodax API
func (c *Client) CallPublic(path string) (*http.Response, error) {
	return http.DefaultClient.Get(c.PublicBaseURL + "/" + path)
}

// CallPrivate generic function to call private Indodax API
func (c *Client) CallPrivate(method Method, params map[string]string) (*http.Response, error) {
	// set params
	ts := time.Now()
	formParams := url.Values{}
	formParams.Set("method", string(method))
	formParams.Set("timestamp", strconv.FormatInt(ts.UnixNano()/int64(time.Millisecond), 10))
	formParams.Set("recvWindow", strconv.FormatInt(ts.Add(c.ReceiveWindow).UnixNano()/int64(time.Millisecond), 10))
	for key, val := range params {
		formParams.Add(key, val)
	}

	// calculate signing
	h := hmac.New(sha512.New, []byte(c.SecretKey))
	h.Write([]byte(formParams.Encode()))
	signingKey := hex.EncodeToString(h.Sum(nil))

	// build request
	req, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewBufferString(formParams.Encode()))
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", c.APIKey)
	req.Header.Set("Sign", signingKey)

	return http.DefaultClient.Do(req)
}
