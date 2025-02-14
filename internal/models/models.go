package models

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory struct {
		Received []ReceivedCoin `json:"received"`
		Sent     []SentCoin     `json:"sent"`
	} `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoin struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}
