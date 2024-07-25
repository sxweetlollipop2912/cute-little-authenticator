package vault

import "little-auth/models"

type Vault interface {
	Get(indexer *models.Indexer) (*models.Secret, error)
	Set(indexer *models.Indexer, secret *models.Secret) error
	Delete(indexer *models.Indexer) error
}

func New() Vault {
	return &vaultImpl{
		verySecretMap: make(map[models.Indexer]*models.Secret),
	}
}

// TODO: don't be lazy
type vaultImpl struct {
	verySecretMap map[models.Indexer]*models.Secret
}

func (v *vaultImpl) Get(indexer *models.Indexer) (*models.Secret, error) {
	if indexer == nil {
		return nil, ErrIndexerNil
	}
	if secret, ok := v.verySecretMap[*indexer]; ok {
		return secret, nil
	}
	return nil, ErrIndexerNotFound
}

func (v *vaultImpl) Set(indexer *models.Indexer, secret *models.Secret) error {
	if indexer == nil {
		return ErrIndexerNil
	}
	if secret == nil {
		return ErrSecretNil
	}
	v.verySecretMap[*indexer] = secret
	return nil
}

func (v *vaultImpl) Delete(indexer *models.Indexer) error {
	if indexer == nil {
		return ErrIndexerNil
	}
	delete(v.verySecretMap, *indexer)
	return nil
}
