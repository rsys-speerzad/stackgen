package configs

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config map[string]string

func ParseEnv(configFilePath string) {
	// open config file
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("error while parsing config file: %s error: %v", configFilePath, err)
	}
	defer file.Close()
	// read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error while parsing config file: %v", err)
	}
	// unmarshal the JSON content into Config map
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("error while parsing config file: %v", err)
	}
	// set environment variables from the config map
	for key, value := range config {
		if err := os.Setenv(key, value); err != nil {
			log.Fatalf("error while parsing config file: %v", err)
		}
	}
}
