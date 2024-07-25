package utils

import "time"

// RFC-6238
func TotpCounterFromNow(period uint64) uint64 {
	return uint64(time.Now().UnixMilli()/1000) / period
}
