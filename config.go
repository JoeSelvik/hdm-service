package main

import (
	"encoding/json"
	"os"
	"time"
)

// A Config represents the application's configuration details.
type Config struct {
	BracketName   string
	StartTime     time.Time
	EndTime       time.Time
	FbGroupId     int `json:"fb_group_id"`
	FbAccessToken string
	FbSession     string
	FbUserId      int
	DbPath        string
	Port          int
}

// NewConfig reads a config from the config file.
//
// The default config file location is $GOPATH/src/github.com/JoeSelvik/hdm-service/config.json, but it may be
// overridden with the BALTOSVC_CONFIG_FILE environment variable. Environment variables in HDMSVC_CONFIG_FILE
// will be expanded, but ~ will not work.
//
// Environment variables in path configs will also have environment variables expanded.
//
// NewConfig will panic if any values are unset.
func NewConfig() *Config {
	// Get the config file path. Defaults to `$GOPATH/src/github.com/JoeSelvik/hdm-service/config.json`
	configFilePath := os.Getenv("HDMSVC_CONFIG_FILE")
	if configFilePath == "" {
		configFilePath = "$GOPATH/src/github.com/JoeSelvik/hdm-service/config.json"
	}

	// Expand environment variables in the config file path
	configFilePath = os.ExpandEnv(configFilePath)

	// Try to open the config file
	f, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Parse the config file
	config := new(Config)
	d := json.NewDecoder(f)
	err = d.Decode(config)
	if err != nil {
		panic(err)
	}

	// Expand environment variables in the path fields

	// Panic over unset, required config variables
	if config.FbGroupId == 0 {
		panic("FbGroupId is not set")
	}

	return config
}
