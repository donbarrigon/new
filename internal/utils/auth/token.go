package auth

import (
	"crypto/rand"
	"donbarrigon/new/internal/utils/logs"
	"encoding/hex"
)

func GenerateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		logs.Error("No se creo el token: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}
