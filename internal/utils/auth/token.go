package auth

import (
	"crypto/rand"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
	"encoding/hex"
	"time"
)

func GenerateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		logs.Error("No se creo el token: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

func expiresAt() time.Time {
	return time.Now().Add(time.Duration(config.SessionLifetime) * time.Minute)
}
