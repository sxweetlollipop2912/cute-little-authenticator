package config

import (
	"little-auth/models"
)

const (
	MAX_FILE_BYTE_LENGTH = 1024
	SERVICE_NAME         = "my-authenticator"
)

var (
	SUPPORTED_ALGO_TYPES = []models.AlgoType{models.TOTP}
)
