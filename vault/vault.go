package vault

import (
	"bytes"
	"encoding/gob"
	"github.com/99designs/keyring"
	zlog "github.com/rs/zerolog/log"
	"little-auth/config"
	"little-auth/crypto"
	"little-auth/crypto/aesgcm"
	"little-auth/models"
	"little-auth/utils"
)

type Vault interface {
	Get(indexer *models.Indexer) (*models.Secret, error)
	Set(indexer *models.Indexer, secret *models.Secret) error
	Delete(indexer *models.Indexer) error
	GetAllIndexes() ([]*models.Indexer, error)
}

func New(otps ...VaultOption) (Vault, error) {
	keyr, err := keyring.Open(keyring.Config{
		ServiceName: config.SERVICE_NAME,
	})
	if err != nil {
		zlog.Err(err).Msg("failed to open keyring")
		return nil, err
	}

	impl := &vaultImpl{
		NameFromIndexer: nameFromIndexer,
		IndexerFromName: indexerFromName,
		GetDir:          getDir,
		newKey:          newKey,
		keyring:         keyr,
		cryptor:         aesgcm.New(aesgcm.AES256),
	}
	for _, opt := range otps {
		opt(impl)
	}
	return impl, nil
}

type vaultImpl struct {
	NameFromIndexer func(*models.Indexer) string
	IndexerFromName func(string) *models.Indexer
	GetDir          func() (string, error)
	newKey          func(int) ([]byte, error)
	keyring         keyring.Keyring
	cryptor         crypto.Cryptor
}

func (v *vaultImpl) Get(indexer *models.Indexer) (*models.Secret, error) {
	var (
		dir         string
		secretModel models.Secret
		err         error
	)
	if indexer == nil {
		return nil, ErrIndexerNil
	}

	// Get directory
	if dir, err = v.GetDir(); err != nil {
		zlog.Err(err).Msg("failed to get directory")
		return nil, err
	}

	// Get secret content
	var secret []byte
	if secret, err = utils.ReadBytesFromFile(dir, v.NameFromIndexer(indexer)); err != nil {
		zlog.Err(err).Msg("failed to read secret from file")
		return nil, err
	}

	// Get key
	var key keyring.Item
	if key, err = v.keyring.Get(v.NameFromIndexer(indexer)); err != nil {
		zlog.Err(err).Msg("failed to get key")
		return nil, err
	}

	// Decrypt secret
	var decrypted []byte
	if decrypted, err = v.cryptor.Decrypt(secret, key.Data); err != nil {
		return nil, err
	}

	// Convert decrypted secret to models.Secret
	if err = gob.NewDecoder(bytes.NewBuffer(decrypted)).
		Decode(&secretModel); err != nil {
		zlog.Err(err).Msg("failed to decode secret")
		return nil, err
	}

	return &secretModel, nil
}

func (v *vaultImpl) Set(indexer *models.Indexer, secret *models.Secret) error {
	var (
		dir         string
		indexerName = v.NameFromIndexer(indexer)
		err         error
	)

	if indexer == nil {
		return ErrIndexerNil
	}
	if secret == nil {
		return ErrSecretNil
	}

	// Encode secret
	var secretBuffer bytes.Buffer
	if err = gob.NewEncoder(&secretBuffer).Encode(secret); err != nil {
		zlog.Err(err).Msg("failed to encode secret")
		return err
	}

	// Encrypt secret
	var vaultKey, encrypted []byte
	if vaultKey, err = v.newKey(v.cryptor.KeyLength()); err != nil {
		zlog.Err(err).Msg("failed to generate key")
		return err

	}
	if encrypted, err = v.cryptor.Encrypt(secretBuffer.Bytes(), vaultKey); err != nil {
		return err
	}

	// Get directory
	if dir, err = v.GetDir(); err != nil {
		zlog.Err(err).Msg("failed to get directory")
		return err
	}

	// Save secret content
	if err = utils.WriteBytesToFile(dir, indexerName, encrypted); err != nil {
		zlog.Err(err).Msg("failed to write secret to file")
		return err
	}

	// Save key
	if err = v.keyring.Set(keyring.Item{
		Key:  indexerName,
		Data: vaultKey,
	}); err != nil {
		zlog.Err(err).Msg("failed to set key")
		return err
	}

	return nil
}

func (v *vaultImpl) Delete(indexer *models.Indexer) error {
	var (
		dir string
		err error
	)

	if indexer == nil {
		return ErrIndexerNil
	}

	// Get directory
	if dir, err = v.GetDir(); err != nil {
		zlog.Err(err).Msg("failed to get directory")
		return err
	}

	if err = utils.DeleteFile(dir, v.NameFromIndexer(indexer)); err != nil {
		zlog.Err(err).Msg("failed to delete secret file")
		return err
	}

	if err = v.keyring.Remove(v.NameFromIndexer(indexer)); err != nil {
		zlog.Err(err).Msg("failed to remove key")
		return err
	}

	return nil
}

func (v *vaultImpl) GetAllIndexes() ([]*models.Indexer, error) {
	var (
		dir       string
		filenames []string
		indexes   []*models.Indexer
		err       error
	)

	// Get directory
	if dir, err = v.GetDir(); err != nil {
		zlog.Err(err).Msg("failed to get directory")
		return nil, err
	}

	if filenames, err = utils.ListFiles(dir); err != nil {
		zlog.Err(err).Msg("failed to list files")
		return nil, err
	}

	// Convert file names to Indexer
	indexes = make([]*models.Indexer, len(filenames))
	for i := range filenames {
		indexes[i] = v.IndexerFromName(filenames[i])
	}
	return indexes, nil
}
