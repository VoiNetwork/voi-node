package utils

import (
	"fmt"
	"log"
	"os"
)

// TODO: Separate network io from blockchain network

const (
	testNet          = "testnet"
	betaNet          = "betanet"
	mainNet          = "mainnet"
	envNetworkVar    = "VOINETWORK_NETWORK"
	envGenesisURLVar = "VOINETWORK_GENESIS"
	envProfileVar    = "VOINETWORK_PROFILE"
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
			StatusURL:   "https://testnet-api.voi.nodly.io/v2/status",
			ArchivalDNS: "voitest.voi.network",
		}, nil
	case betaNet:
		return Network{
			Name:        betaNet,
			StatusURL:   "https://betanet-api.voi.nodly.io/v2/status",
			ArchivalDNS: "betanet-voi.network",
		}, nil
	case mainNet:
		return Network{
			Name:        mainNet,
			StatusURL:   "https://mainnet-api.voi.nodly.io/v2/status",
			ArchivalDNS: "mainnet-voi.network",
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
