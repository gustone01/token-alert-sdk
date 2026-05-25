package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func dedupKey(event Event) string {
	tokenHash := hashToken(event.TokenHint)
	return fmt.Sprintf("tokenalert:dedup:%s:%s:%d", event.Platform, tokenHash, event.Code)
}

func hashToken(token string) string {
	if token == "" {
		return "empty"
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:8])
}
