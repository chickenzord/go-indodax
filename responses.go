package indodax

import (
	"fmt"
	"strconv"
	"time"
)

// ResponseMeta response metadata
type ResponseMeta struct {
	Success      int    `json:"success"`
	ErrorMessage string `json:"error"`
	ErrorCode    string `json:"error_code"`
}

// Error return nil if success
func (r *ResponseMeta) Error() error {
	if r.Success == 1 {
		return nil
	}

	return fmt.Errorf(r.ErrorMessage)
}

// Info user info
type Info struct {
	UserID            string  `json:"user_id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	ServerTimeSeconds int64   `json:"server_time"`
	BalanceAvailable  Balance `json:"balance"`
	BalanceHold       Balance `json:"balance_hold"`
}

// ServerTime returns server time
func (i *Info) ServerTime() time.Time {
	return time.Unix(0, i.ServerTimeSeconds*int64(time.Second))
}

// Balance contains crypto balance
type Balance map[string]interface{}

// Get balance amount by currency
func (b Balance) Get(currency string) float64 {
	switch val := b[currency].(type) {
	case float64:
		return val
	case string:
		balance, err := strconv.ParseFloat(val, 64)
		if err != nil {
			panic(err)
		}
		return balance
	default:
		panic(fmt.Errorf("Cannot parse value of %q", val))
	}
}

// NonEmpty returns new Balance with empty amounts filtered out
func (b Balance) NonEmpty() Balance {
	result := Balance{}
	for key := range b {
		if balance := b.Get(key); balance > 0 {
			result[key] = balance
		}
	}
	return result
}

// GetInfoResponse response from getInfo method
type GetInfoResponse struct {
	ResponseMeta `json:",inline"`
	Info         Info `json:"return"`
}
