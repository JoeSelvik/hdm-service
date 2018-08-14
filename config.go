package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	// fb created_time str:      2017-03-04T13:05:20+0000
	// sqlite CURRENT_TIMESTAMP: 2017-03-06 15:36:17
	// Golang template time      Mon, 01/02/06, 03:04PM
	// HDM golang template       Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)
	GoTimeLayout = "2006-01-02T15:04:05+0000"
)

// A Configuration represents the application's configuration details.
type Configuration struct {
	BracketName     string
	StartTime       time.Time
	StartTimeString string `json:"start_time"`
	EndTime         time.Time
	EndTimeString   string `json:"end_time"`
	FbGroupId       int    `json:"fb_group_id"`
	FbAccessToken   string `json:"fb_access_token"`
	FbSession       string
	FbUserId        int
	DbPath          string `json:"db_path"`
	DbSetupScript   string `json:"db_setup_script"` // optional
	DbTestPath      string `json:"db_test_path"`    // optional
	Port            int
}

// NewConfig reads a configuration values from config.json.
//
// The default config file location is in the project root dir, but it may be overridden with the HDMSVC_CONFIG_FILE
// environment variable. Environment variables in HDMSVC_CONFIG_FILE will be expanded, but ~ will not work. Environment
// variables in path configs will also have environment variables expanded.
//
// NewConfig will panic if any configuration values are unset.
//
// Note - Probably better to use a file format that supports comments, like yaml
func NewConfig() *Configuration {
	// get the config file
	configFile := os.Getenv("HDMSVC_CONFIG_FILE")
	if configFile == "" {
		configFile = "config.json"
	}

	// expand environment variables in the config file path
	configFile = os.ExpandEnv(configFile)

	// open the config file
	f, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// parse the config file
	config := new(Configuration)
	d := json.NewDecoder(f)
	err = d.Decode(config)
	if err != nil {
		panic(err)
	}

	// todo: If given, expand environment variables in the path fields

	// panic over unset required config variables
	if config.FbGroupId == 0 {
		panic("fb_group_id is not set")
	}
	if config.FbAccessToken == "" {
		panic("fb_access_token is not set")
	}
	if config.DbPath == "" {
		panic("db_path is not set")
	}

	// todo: Is this the right place to handle start and end time, with two values?
	if config.StartTimeString == "" {
		panic("start_time is not set")
	}
	t, err := time.Parse(GoTimeLayout, config.StartTimeString)
	if err != nil {
		msg := fmt.Sprintf("Could not parse given start_time, check formatting. Try: %s", GoTimeLayout)
		panic(msg)
	}
	config.StartTime = t

	if config.EndTimeString == "" {
		panic("end_time is not set")
	}
	t, err = time.Parse(GoTimeLayout, config.EndTimeString)
	if err != nil {
		msg := fmt.Sprintf("Could not parse given end_time, check formatting. Try: %s", GoTimeLayout)
		panic(msg)
	}
	config.EndTime = t

	return config
}
