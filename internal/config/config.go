package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL	  string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/workspace/github.com/nunseik/gator/.gatorconfig.json"

func Read () (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg Config) SetUser() error{
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFileName, nil
}
