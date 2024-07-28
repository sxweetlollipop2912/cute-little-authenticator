package vault

import (
	"crypto/rand"
	"fmt"
	"little-auth/config"
	"little-auth/models"
	"little-auth/utils"
	"os"
	"strings"
)

const (
	PATH_ISSUER_SEPARATOR = "___"
)

func nameFromIndexer(indexer *models.Indexer) string {
	if indexer == nil {
		return ""
	}
	var name string
	if indexer.Path != "" && indexer.Issuer != "" {
		name = fmt.Sprintf("%s%s%s", indexer.Path, PATH_ISSUER_SEPARATOR, indexer.Issuer)
	} else if indexer.Path == "" && indexer.Issuer != "" {
		name = indexer.Issuer
	} else if indexer.Path != "" && indexer.Issuer == "" {
		name = indexer.Path
	}
	return utils.NormalizeFileName(name)
}

func indexerFromName(name string) *models.Indexer {
	if name == "" {
		return nil
	}
	indexer := &models.Indexer{}
	if i := strings.Index(name, PATH_ISSUER_SEPARATOR); i != -1 {
		indexer.Path = name[:i]
		indexer.Issuer = name[i+len(PATH_ISSUER_SEPARATOR):]
	} else {
		indexer.Path = name
	}
	return indexer
}

func newKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func getDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homedir + "/Library/Application Support/" + config.SERVICE_NAME + "/", nil
}
