package keyman

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type keyMapping struct {
	IpAddress string
	ApiKey    string
}

type KeyManager struct {
	keyMappings []keyMapping
}

func NewKeyManager() KeyManager {
	return KeyManager{}
}

func GetDefaultKeyStoreFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to locate user home directory.")
	}

	keyStoreFilePath := filepath.Join(homeDir, ".huego", "keys.json")

	return keyStoreFilePath
}

type KeyStoreNotExistsError struct {
	keyStoreFilePath string
}

func (e *KeyStoreNotExistsError) Error() string {
	return fmt.Sprintf("can't load key store file: %s", e.keyStoreFilePath)
}

func checkKeyStoreFileExists(keyStoreFilePath string) error {
	if _, err := os.Stat(keyStoreFilePath); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			return &KeyStoreNotExistsError{keyStoreFilePath}
		} else {
			// some other unexpected error has occurred
			return err
		}
	}

	return nil
}

type keyStore struct {
	KeyMappings []keyMapping
}

func (k *KeyManager) LoadFromKeyStore(keyStoreFilePath string) error {
	if err := checkKeyStoreFileExists(keyStoreFilePath); err != nil {
		return err
	}

	// key store file exists
	bytes, err := os.ReadFile(keyStoreFilePath)
	if err != nil {
		return errors.New("failed to read key store file contents")
	}

	var keyStore keyStore
	if err := json.Unmarshal(bytes, &keyStore); err != nil {
		return errors.New("failed to unmarshal json from key store file")
	}

	k.keyMappings = keyStore.KeyMappings

	return nil
}

func (c KeyManager) SaveToKeyStore(keyStoreFilePath string) error {
	parentDir := filepath.Dir(keyStoreFilePath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		// parent key store dir does not exist; create it
		err := os.MkdirAll(parentDir, 0770)
		if err != nil {
			return errors.New("failed to create parent directory to save key store file")
		}
	} else if err != nil {
		// some other unexpected error has occurred
		return fmt.Errorf("unexpected error: %s", err.Error())
	}

	// parent key store dir exists; write key store file
	keyStore := keyStore{KeyMappings: c.keyMappings}
	bytes, err := json.Marshal(keyStore)
	if err != nil {
		return errors.New("failed to marshal key store")
	}
	err = os.WriteFile(keyStoreFilePath, bytes, 0770)
	if err != nil {
		return errors.New("failed to write key store file")
	}

	return nil
}

func (c *KeyManager) SetApiKey(ip string, apiKey string) error {
	// check for and update existing API key entry
	for i := 0; i < len(c.keyMappings); i++ {
		hub := &c.keyMappings[i]
		if hub.IpAddress == ip {
			hub.ApiKey = apiKey
			return nil
		}
	}

	// otherwise add new one
	c.keyMappings = append(c.keyMappings, keyMapping{
		IpAddress: ip,
		ApiKey:    apiKey,
	})
	return nil
}

func (c KeyManager) GetApiKey(ip string) (string, error) {
	for _, savedHub := range c.keyMappings {
		if savedHub.IpAddress == ip {
			return savedHub.ApiKey, nil
		}
	}

	return "", nil
}
