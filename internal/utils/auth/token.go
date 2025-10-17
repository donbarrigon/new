package auth

import (
	"crypto/rand"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/logs"
	"encoding/hex"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, e := rand.Read(bytes); e != nil {
		logs.Error("Error al generar el token")
		return "", err.New(err.INTERNAL, "Error al generar el token", e)
	}
	return hex.EncodeToString(bytes), nil
}
