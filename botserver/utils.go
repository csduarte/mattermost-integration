package botserver

import (
	"encoding/json"
	"os"
)

// parseConfig unmarshals json data for Config
func parseConfig(path string) (Config, error) {
	config := Config{}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return config, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
