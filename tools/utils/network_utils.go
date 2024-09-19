package utils

import (
	"fmt"
	"log"
	"os"
)

// TODO: Separate network io from blockchain network

const (
	testNet            = "testnet"
	betaNet            = "betanet"
	mainNet            = "mainnet"
	envNetworkVar      = "VOINETWORK_NETWORK"
	envGenesisURLVar   = "VOINETWORK_GENESIS"
	envProfileVar      = "VOINETWORK_PROFILE"
	envTelemetryVar    = "VOINETWORK_TELEMETRY_NAME"
	envOverwriteConfig = "VOINETWORK_OVERWRITE_CONFIG"
)

type Network struct {
	Name        string
	StatusURL   string
	ArchivalDNS string
}

type NetworkUtils struct{}

func (nu NetworkUtils) NewNetwork(name string) (Network, error) {
	switch name {
	case testNet:
		return Network{
			Name:        testNet,
			StatusURL:   "https://testnet-api.voi.nodely.io/v2/status",
			ArchivalDNS: "voitest.testnet-voi.network",
		}, nil
	case betaNet:
		return Network{
			Name:        betaNet,
			StatusURL:   "https://betanet-api.voi.nodely.io/v2/status",
			ArchivalDNS: "voibeta.betanet-voi.network",
		}, nil
	case mainNet:
		return Network{
			Name:        mainNet,
			StatusURL:   "https://mainnet-api.voi.nodely.dev/v2/status",
			ArchivalDNS: "voimain.mainnet-voi.network",
		}, nil
	default:
		return Network{}, fmt.Errorf("unsupported network: %s", name)
	}
}

func (nu NetworkUtils) CheckIfPredefinedNetwork(network string) bool {
	switch network {
	case testNet:
		return true
	case betaNet:
		return true
	case mainNet:
		return true
	default:
		return false
	}
}

func (nu NetworkUtils) GetEnvNetworkVar() string {
	return envNetworkVar
}

func (nu NetworkUtils) GetNetworkFromEnv() (string, bool) {
	network := os.Getenv(envNetworkVar)
	if network != "" {
		log.Printf("Using network from environment variable: %s", network)
		return network, true
	}
	return "", false
}

// TODO: Remove duplication by extracting the common code to a function
func (nu NetworkUtils) GetEnvProfileVar() string {
	return envProfileVar
}

func (nu NetworkUtils) GetProfileFromEnv() (string, bool) {
	profile := os.Getenv(envProfileVar)
	if profile != "" {
		log.Printf("Using profile from environment variable: %s", profile)
		return profile, true
	}
	return "", false
}

func (nu NetworkUtils) GetGenesisFromEnv() (string, bool) {
	genesisURL := os.Getenv(envGenesisURLVar)
	if genesisURL != "" {
		return genesisURL, true
	}
	return "", false
}

func (nu NetworkUtils) GetEnvTelemetryVar() string {
	return envTelemetryVar
}

func (nu NetworkUtils) GetTelemetryNameFromEnv() (string, bool) {
	telemetryName := os.Getenv(envTelemetryVar)
	legacyTelemetryName := os.Getenv("TELEMETRY_NAME")

	if telemetryName != "" {
		return telemetryName, true
	}
	if legacyTelemetryName != "" {
		return legacyTelemetryName, true
	}
	return "", false
}

func (nu NetworkUtils) GetEnvOverwriteConfig() string {
	return envOverwriteConfig
}

func (nu NetworkUtils) GetOverwriteConfigFromEnv() (bool, bool) {
	overwriteConfig := os.Getenv(envOverwriteConfig)
	if overwriteConfig == "true" {
		return true, true
	}
	if overwriteConfig == "false" {
		return false, true
	}
	return false, false
}
