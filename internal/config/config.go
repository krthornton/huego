package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Hub struct {
	IpAddress string
	ApiKey    string
}

type Configuration struct {
	Hubs []Hub
}

func NewConfiguration() Configuration {
	return Configuration{}
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to locate user home directory.")
	}

	configFilePath := filepath.Join(homeDir, ".huego", "config.json")

	return configFilePath
}

type ConfigFileNotExists struct {
	configFilePath string
}

func (e *ConfigFileNotExists) Error() string {
	return fmt.Sprintf("can't load config file: %s", e.configFilePath)
}

func LoadConfiguration() (Configuration, error) {
	var config Configuration
	configFilePath := getConfigFilePath()
	if _, err := os.Stat(configFilePath); err == nil {
		// config file exists
		bytes, err := os.ReadFile(configFilePath)
		if err != nil {
			panic("Failed to read config file contents")
		}

		if err := json.Unmarshal(bytes, &config); err != nil {
			panic("Failed to unmarshal json from config file")
		}

		return config, nil
	} else if os.IsNotExist(err) {
		// file does not exist
		return config, &ConfigFileNotExists{configFilePath}
	} else {
		// some other unexpected error has occurred
		panic(fmt.Sprintf("Unexpected error: %s", err.Error()))
	}
}

func (c Configuration) SaveConfiguration() {
	configFilePath := getConfigFilePath()
	parentDir := filepath.Dir(configFilePath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		// parent config dir does not exist; create it
		err := os.MkdirAll(parentDir, 0770)
		if err != nil {
			panic("Failed to create parent directory to store config file")
		}
	} else if err != nil {
		// some other unexpected error has occurred
		panic(fmt.Sprintf("Unexpected error: %s", err.Error()))
	}

	// parent config dir exists; write config file
	bytes, err := json.Marshal(c)
	if err != nil {
		panic("Failed to marshal configuration")
	}
	err = os.WriteFile(configFilePath, bytes, 0770)
	if err != nil {
		panic("Failed to write configuration file")
	}
}

func (c *Configuration) SetApiKey(ip string, apiKey string) error {
	// check for and update existing API key entry
	for i := 0; i < len(c.Hubs); i++ {
		hub := &c.Hubs[i]
		if hub.IpAddress == ip {
			hub.ApiKey = apiKey
			return nil
		}
	}

	// otherwise add new one
	c.Hubs = append(c.Hubs, Hub{
		IpAddress: ip,
		ApiKey:    apiKey,
	})
	return nil
}

func (c Configuration) GetApiKey(ip string) (string, error) {
	for _, savedHub := range c.Hubs {
		if savedHub.IpAddress == ip {
			return savedHub.ApiKey, nil
		}
	}

	return "", nil
}
