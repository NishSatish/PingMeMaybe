package utils

import (
	"time"
)

func IsOlderThanOneDay(createdAt time.Time) bool {
	return time.Since(createdAt) > 24*time.Hour
}
