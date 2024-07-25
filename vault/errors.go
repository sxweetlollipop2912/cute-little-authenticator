package vault

import "errors"

var (
	ErrIndexerNotFound = errors.New("indexer not found")
	ErrIndexerNil      = errors.New("indexer is nil")
	ErrSecretNil       = errors.New("secret is nil")
)
