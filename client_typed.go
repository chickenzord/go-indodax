package indodax

import "encoding/json"

// TypedClient Indodax client with typed responses
type TypedClient struct {
	Client
}

// NewDefaultTypedClient create new default Indodax client with typed resources
func NewDefaultTypedClient(apiKey, secretKey string) *TypedClient {
	return &TypedClient{
		Client: *NewDefaultClient(apiKey, secretKey),
	}
}

// GetInfo This method gives user balances and server's timestamp.
func (c *TypedClient) GetInfo() (*GetInfoResponse, error) {
	r, err := c.CallPrivate("getInfo", map[string]string{})
	if err != nil {
		return nil, err
	}

	var result GetInfoResponse
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}
	if err := result.Error(); err != nil {
		return nil, err
	}

	return &result, nil
}
