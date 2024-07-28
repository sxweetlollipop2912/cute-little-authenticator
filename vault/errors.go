package vault

import "errors"

var (
	ErrIndexerNil = errors.New("indexer is nil")
	ErrSecretNil  = errors.New("secret is nil")
)
