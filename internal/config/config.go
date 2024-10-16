// Responsible for reading and writing the JSON file
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// filename
const configFileName = ".gatorconfig.json"

// Export Config struct that represents the JSON file structure, including struct tags.
type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Export a Read function that reads the JSON file and returns a Config struct.
func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	// Open the JSON file
	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening json file: %s", err)
	}
	defer file.Close()

	// Read file 's json
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Export SetUser method that setting username on struct
func (cfg *Config) SetUser(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("there is no username provided")
	}
	// Set the new username
	cfg.CurrentUserName = username

	return write(*cfg)

}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(homeDir, configFileName)
	// fullPath := "/mnt/c/Users/janis/python_projects/boot_dev_project" + "/" + configFileName
	return fullPath, nil
}

func write(cfg Config) error {

	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	//Write the jsondata into file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(cfg)
	if err != nil {
		return fmt.Errorf("error in encoding struct into json file:%s", err)
	}

	return nil
}
