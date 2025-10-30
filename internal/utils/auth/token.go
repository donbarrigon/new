package auth

import (
	"crypto/rand"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/logs"
	"encoding/hex"
	"math/big"
)

// genera un token exadecimal de 32 bytes
func GenerateHexToken() (string, error) {
	bytes := make([]byte, 32)
	if _, e := rand.Read(bytes); e != nil {
		logs.Error("Error al generar el token")
		return "", err.New(err.INTERNAL, "Error al generar el token", e)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateVerificationCode() (string, error) {
	code := ""
	for range 6 {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += num.String()
	}
	return code, nil
}
