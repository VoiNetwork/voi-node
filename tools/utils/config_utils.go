package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	genesisJSONPathFmt             = "/algod/configuration/%s/genesis.json"
	configJSONPathFmt              = "/algod/configuration/%s/%s/config.json"
	envIncomingConnectionsLimitVar = "VOINETWORK_INCOMING_CONNECTIONS_LIMIT"
)

type ConfigUtils struct{}

func (cu ConfigUtils) HandleConfiguration(urlSet bool, genesisURL string, network string, profile string, overwriteConfig bool, algodDataDir string) {
	fu := FileUtils{}
	nu := NetworkUtils{}

	if urlSet {
		log.Printf("Using genesis and configuration URLs from environment variables: %s", genesisURL)
		if err := nu.DownloadNetworkConfiguration(genesisURL, algodDataDir); err != nil {
			fmt.Printf("Failed to download network configuration: %v", err)
			os.Exit(1)
		}
	} else {
		if err := fu.CopyGenesisConfigurationFromFilesystem(network, overwriteConfig, genesisJSONPathFmt, algodDataDir); err != nil {
			fmt.Printf("Failed to copy network configuration: %v", err)
			os.Exit(1)
		}
	}

	if err := fu.CopyAlgodConfigurationFromFilesystem(network, profile, overwriteConfig, configJSONPathFmt, algodDataDir); err != nil {
		fmt.Printf("Failed to copy network configuration: %v", err)
		os.Exit(1)
	}

	configPath := filepath.Join(algodDataDir, "config.json")
	err := overrideConfigurationVariable(configPath, "IncomingConnectionsLimit", os.Getenv(envIncomingConnectionsLimitVar))
	if err != nil {
		return
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	log.Print(string(content))
}

func overrideConfigurationVariable(configPath, key, value string) error {
	if value == "" {
		return nil
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %v", err)
	}

	if value == "true" || value == "false" {
		boolValue := value == "true"
		config[key] = boolValue
		log.Printf("Overwriting boolean value for %s: %t", key, boolValue)
	} else {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			config[key] = intValue
			log.Printf("Overwriting integer value for %s: %d", key, intValue)
		} else {
			config[key] = value
			log.Printf("Overwriting string value for %s: %s", key, value)
		}
	}

	modifiedConfig, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal modified config: %v", err)
	}

	if err := os.WriteFile(configPath, modifiedConfig, 0644); err != nil {
		return fmt.Errorf("failed to write modified config back to file: %v", err)
	}

	return nil
}
