package models

import "time"

// Structs to represent the JSON data
type OnlineResponse struct {
	ContractMap map[string]bool  `json:"contractMap"`
	TokenInfo   OnlineTokenInfo  `json:"tokenInfo"`
	PageSize    int              `json:"page_size"`
	Code        int              `json:"code"`
	Data        []OnlineDataItem `json:"data"`
}

type OnlineTokenInfo struct {
	TokenId      string `json:"tokenId"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenName    string `json:"tokenName"`
	TokenDecimal int    `json:"tokenDecimal"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
	TokenLogo    string `json:"tokenLogo"`
	TokenLevel   string `json:"tokenLevel"`
	IssuerAddr   string `json:"issuerAddr"`
	VIP          bool   `json:"vip"`
}

type OnlineDataItem struct {
	Amount          string `json:"amount"`
	Status          int    `json:"status"`
	ApprovalAmount  string `json:"approval_amount"`
	BlockTimestamp  int64  `json:"block_timestamp"`
	Block           int    `json:"block"`
	From            string `json:"from"`
	To              string `json:"to"`
	Hash            string `json:"hash"`
	ContractAddress string `json:"contract_address"`
	Confirmed       int    `json:"confirmed"`
	ContractType    string `json:"contract_type"`
	ContractTypeInt int    `json:"contractType"`
	Revert          int    `json:"revert"`
	ContractRet     string `json:"contract_ret"`
	FinalResult     string `json:"final_result"`
	EventType       string `json:"event_type"`
	IssueAddress    string `json:"issue_address"`
	Decimals        int    `json:"decimals"`
	TokenName       string `json:"token_name"`
	ID              string `json:"id"`
	Direction       int    `json:"direction"`
}

// Define a struct that matches the JSON structure
type OnlineTokenData struct {
	TokenId      string `json:"tokenId"`
	TokenName    string `json:"tokenName"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
}

// Struct for transactions
type TransactionUpdateOnline struct {
	//gorm.Model
	Id                 uint      `gorm:"primaryKey"`
	Receivedamount     float64   `json:"receivedamount,omitempty"`
	Receivedcurrency   string    `json:"receivedcurrency,omitempty"`
	Status             string    `json:"status,omitempty"`
	Response_hash      string    `json:"response_hash,omitempty"`
	Response_from      string    `json:"response_from,omitempty"`
	Response_to        string    `json:"response_to,omitempty"`
	Response_timestamp time.Time `json:"response_timestamp,omitempty"`
	Response_json      string    `json:"response_json,omitempty"`
}

type PayRequest struct {
	Cid            string `json:"cid" form:"cid"`
	Price_currency string `json:"price_currency" form:"price_currency"`
	Price_amount   string `json:"price_amount" form:"price_amount"`
	Sender_name    string `json:"sender_name" form:"sender_name"`
	Sender_email   string `json:"sender_email" form:"sender_email"`
	Client_id      uint   `json:"client_id" form:"client_id"`
	Pay_type       int    `json:"pay_type" form:"pay_type"`
	Crypto_id      int    `json:"crypto_id" form:"crypto_id"`
	Customerrefid  string `json:"customerrefid" form:"customerrefid"`
}
type PayResponse struct {
	Qr_code     string  `json:"qr_code"`
	Address     string  `json:"address"`
	Amount      float64 `json:"amount"`
	Transid     string  `json:"transid"`
	Coinicon    string  `json:"coinicon"`
	Coinnetwork string  `json:"coinnetwork"`
	Coin_id     int     `json:"coin_id"`
}
