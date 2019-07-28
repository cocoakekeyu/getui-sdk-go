package utils

import (
	"strconv"
	"time"
)

func GenerateRequestID() string {
	RequestID := strconv.FormatInt(time.Now().UnixNano(), 12)
	return RequestID
}
