package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: Separate network io from blockchain network

const (
	// TODO: Create an enum to hold predefined network values
	testNet          = "testnet"
	envNetworkVar    = "VOINETWORK_NETWORK"
	envGenesisURLVar = "VOINETWORK_GENESIS"
	envProfileVar    = "VOINETWORK_PROFILE"
)

type NetworkUtils struct{}

func (nu NetworkUtils) GetStatusURL(network string) (string, error) {
	switch network {
	case testNet:
		return "https://testnet-api.voi.nodly.io/v2/status", nil
	default:
		return "", fmt.Errorf("unsupported network: %s", network)
	}
}

func (nu NetworkUtils) CheckIfPredefinedNetwork(network string) bool {
	switch network {
	case testNet:
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

func (nu NetworkUtils) DownloadNetworkConfiguration(genesisURL, algodDataDir string) error {
	if err := downloadFile(genesisURL, filepath.Join(algodDataDir, "genesis.json")); err != nil {
		return fmt.Errorf("failed to download genesis.json: %w", err)
	}

	return nil
}

func downloadFile(url, destFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	fu := FileUtils{}
	err = fu.EnsureDirExists(filepath.Dir(destFile))
	if err != nil {
		return err
	}

	out, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", destFile, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing data to %s: %w", destFile, err)
	}

	log.Printf("Successfully downloaded %s to %s", url, destFile)
	return nil
}
