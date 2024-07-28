package vault

import (
	"github.com/99designs/keyring"
	"little-auth/crypto"
	"little-auth/models"
)

type VaultOption func(*vaultImpl)

func WithNameFromIndexerFn(fn func(*models.Indexer) string) VaultOption {
	return func(v *vaultImpl) {
		v.NameFromIndexer = fn
	}
}

func WithIndexerFromNameFn(fn func(string) *models.Indexer) VaultOption {
	return func(v *vaultImpl) {
		v.IndexerFromName = fn
	}
}

func WithNewKeyFn(fn func(int) ([]byte, error)) VaultOption {
	return func(v *vaultImpl) {
		v.newKey = fn
	}
}

func WithKeyring(k keyring.Keyring) VaultOption {
	return func(v *vaultImpl) {
		v.keyring = k
	}
}

func WithCryptor(c crypto.Cryptor) VaultOption {
	return func(v *vaultImpl) {
		v.cryptor = c
	}
}

func WithGetDirFn(fn func() (string, error)) VaultOption {
	return func(v *vaultImpl) {
		v.GetDir = fn
	}
}
