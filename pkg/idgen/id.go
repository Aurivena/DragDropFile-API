package idgen

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

const lenCode = 12

func GenerateID() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenCode)
	newInt := big.NewInt(int64(len(letters)))
	for i := range code {
		num, err := rand.Int(rand.Reader, newInt)
		if err != nil {
			logrus.Errorf("не удалось сгенерировать часть ID: %s", err)
			return "", fmt.Errorf("не удалось сгенерировать часть ID: %w", err)
		}
		code[i] = letters[num.Int64()]
	}

	return string(code), nil
}
