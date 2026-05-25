package core

import "time"

type TokenType string

const (
	TokenTypeAccessToken TokenType = "access_token"
	TokenTypeUserToken   TokenType = "user_token"
)

type EnrichInfo struct {
	AccountID string
	Name      string
	Channel   string
}

type EnrichFunc func(accessToken string) EnrichInfo

type Event struct {
	Platform   Platform
	TokenType  TokenType
	Service    string
	APIPath    string
	Code       int64
	Message    string
	TokenHint  string
	Enrich     EnrichInfo
	OccurredAt time.Time
}
