package indodax

import "encoding/json"

// Ticker crypto ticker values
type Ticker struct {
	High float64 `json:"high,string"`
	Low  float64 `json:"low,string"`
	Last float64 `json:"last,string"`
	Buy  float64 `json:"buy,string"`
	Sell float64 `json:"sell,string"`
}

// TickerResponse reponse from get ticker
type TickerResponse struct {
	Ticker Ticker `json:"ticker"`
}

func (c *TypedClient) GetTicker(pairID string) (*TickerResponse, error) {
	r, err := c.CallPublic("ticker/" + pairID)
	if err != nil {
		return nil, err
	}

	var result TickerResponse
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
