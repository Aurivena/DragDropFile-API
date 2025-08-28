package fileops

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func DecodeFile(fileBase64 string) ([]byte, error) {
	if !strings.HasPrefix(fileBase64, "data:") {
		logrus.Error("invalid base64 format: missing data prefix")
		return nil, fmt.Errorf("invalid base64 format: missing data prefix")
	}

	parts := strings.SplitN(fileBase64, ";base64,", 2)
	if len(parts) != 2 {
		logrus.Error("invalid base64 format: missing base64 separator")
		return nil, fmt.Errorf("invalid base64 format: missing base64 separator")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		logrus.Error("failed to decode base64: %w", err)
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return data, nil
}
