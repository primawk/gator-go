package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string
	Current_user_name string
}

const configFileName = ".gatorconfig.json"

func write(cfg Config) error {
	dir, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(dir, data, 0644)
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
	return homeDir + "/" + configFileName, nil
}

func Read() Config {
	if dir, err := getConfigFilePath(); err != nil {
		fmt.Println(err)
		return Config{}
	} else {
		if data, err := os.ReadFile(dir); err != nil {
			fmt.Println("Error reading file: %v", err)
			return Config{}
		} else {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err != nil {
				fmt.Println(err)
			}
			return cfg
		}
	}
}

func (c Config) SetUser(username string) {
	c.Current_user_name = username
	if err := write(c); err != nil {
		fmt.Println(err)
	}
}
