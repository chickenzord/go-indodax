package indodax

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

type TradeType string

var (
	TradeTypeBuy  TradeType = "buy"
	TradeTypeSell TradeType = "sell"
)

type Trade map[string]string

func (t Trade) Fee() float64 {
	val, err := strconv.ParseFloat(t["fee"], 64)
	if err != nil {
		panic(fmt.Errorf("Cannot get fee for Trade: %+v", t))
	}

	return val
}

func (t Trade) TargetCurrency() string {
	return t["currency"]
}

func (t Trade) TargetAmount() float64 {
	val, err := strconv.ParseFloat(t[t.TargetCurrency()], 64)
	if err != nil {
		panic(fmt.Errorf("Cannot get target amount for Trade: %+v", t))
	}

	return val
}

func (t Trade) BasePrice() float64 {
	val, err := strconv.ParseFloat(t["price"], 64)
	if err != nil {
		panic(fmt.Errorf("Cannot get base price for Trade: %+v", t))
	}

	return val
}

func (t Trade) BaseCurrency() string {
	frags := strings.Split(t.Pair(), "_")
	return frags[len(frags)-1]
}

func (t Trade) Type() string {
	return t["type"]
}

func (t Trade) Pair() string {
	return t["pair"]
}

func (t Trade) TotalCost() float64 {
	// TODO: calculate fee
	return t.BasePrice() * t.TargetAmount()
}

func (t Trade) Time() time.Time {
	val, err := strconv.ParseInt(t["trade_time"], 10, 64)
	if err != nil {
		panic(fmt.Errorf("Cannot get time for Trade: %+v", t))
	}

	return time.Unix(val, 0)
}

func (t Trade) Description() string {
	return fmt.Sprintf("[%s] %f %s for %f %s with a total of %f %s on [%s]",
		strings.ToUpper(t.Type()),
		t.TargetAmount(), strings.ToUpper(t.TargetCurrency()),
		t.BasePrice(), strings.ToUpper(t.BaseCurrency()),
		t.TotalCost(), strings.ToUpper(t.BaseCurrency()),
		t.Time(),
	)
}

type TradeHistory struct {
	Trades []Trade `json:"trades"`
}

type TradeHistoryResponse struct {
	ResponseMeta `json:",inline"`
	TradeHistory TradeHistory `json:"return"`
}

func (c *TypedClient) GetTradeHistory(pair string) (*TradeHistoryResponse, error) {
	r, err := c.CallPrivate(MethodTradeHistory, map[string]string{
		"pair":  pair,
		"order": "desc",
	})
	if err != nil {
		return nil, err
	}

	var result TradeHistoryResponse
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}
	if err := result.Error(); err != nil {
		return nil, err
	}

	return &result, nil
}
