package main

import (
	"flag"
	"fmt"
	"github.com/voinetwork/voi-node/tools/utils"
	"log"
	"os"
	"time"
)

const (
	envCatchupVar      = "VOINETWORK_CATCHUP"
	genesisJSONPathFmt = "/algod/configuration/%s/genesis.json"
	configJSONPathFmt  = "/algod/configuration/%s/%s/config.json"
	algodDataDir       = "/algod/data"
	algodCmd           = "/node/bin/algod"
	catchupCmd         = "/node/bin/catch-catchpoint"
)

var network string
var profile string
var overwriteConfig bool

func init() {
	flag.StringVar(&network, "network", "testnet", "Specify the network (testnet)")
	flag.StringVar(&profile, "profile", "relay", "Specify the profile (archiver, relay)")
	flag.BoolVar(&overwriteConfig, "overwrite-config", true, "Specify whether to overwrite the configuration files (true, false)")
}

func handleConfiguration(urlSet bool, genesisURL, network, profile string, overwriteConfig bool, genesisJSONPathFmt, configJSONPathFmt, algodDataDir string) {
	fu := utils.FileUtils{}
	nu := utils.NetworkUtils{}

	if urlSet {
		log.Printf("Using genesis and configuration URLs from environment variables: %s", genesisURL)
		if err := nu.DownloadNetworkConfiguration(genesisURL, algodDataDir); err != nil {
			fmt.Printf("Failed to download network configuration: %v", err)
			os.Exit(1)
		}
	} else {
		if err := fu.CopyGenesisConfigurationFromFilesystem(network, profile, overwriteConfig, genesisJSONPathFmt, algodDataDir); err != nil {
			fmt.Printf("Failed to copy network configuration: %v", err)
			os.Exit(1)
		}
	}

	if err := fu.CopyAlgodConfigurationFromFilesystem(network, profile, overwriteConfig, configJSONPathFmt, algodDataDir); err != nil {
		fmt.Printf("Failed to copy network configuration: %v", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// use default network if not set as env variable. For backwards compatibility with existing nodes we will default
	// to testnet instead of requiring a network to be set.
	nu := utils.NetworkUtils{}
	envNetwork, networkSet := nu.GetNetworkFromEnv()
	if networkSet {
		network = envNetwork
	}

	genesisURL, urlSet := nu.GetGenesisFromEnv()

	envProfile, profileSet := nu.GetProfileFromEnv()
	if profileSet {
		profile = envProfile
	}

	log.Printf("Network: %s", network)
	log.Printf("Profile: %s", profile)
	log.Printf("Overwrite Config: %t", overwriteConfig)

	handleConfiguration(urlSet, genesisURL, network, profile, overwriteConfig, genesisJSONPathFmt, configJSONPathFmt, algodDataDir)

	pu := utils.ProcessUtils{}

	// Start algod
	done := pu.StartProcess(algodCmd, "-d", algodDataDir)

	// Wait for algod to start
	time.Sleep(5 * time.Second)

	envVar := os.Getenv(envCatchupVar)
	if envVar != "0" && !urlSet {
		// Execute catch-catchpoint
		_, err := pu.ExecuteCommand(catchupCmd)
		if err != nil {
			return
		}
	} else {
		log.Printf("Skipping catchup execution as %s is set to 0", envCatchupVar)
	}

	<-done
}
