package main

import (
	"flag"
	"fmt"
	"github.com/voinetwork/docker-relay-node/tools/utils"
	"log"
	"net"
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

func main() {
	flag.Parse()

	// use default network if not set as env variable. For backwards compatibility with existing nodes we will default
	// to testnet instead of requiring a network to be set.
	nu := utils.NetworkUtils{}
	envNetwork, networkSet := nu.GetNetworkFromEnv()
	if networkSet {
		network = envNetwork
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(err)
			continue
		}

		for _, addr := range addrs {
			fmt.Printf("Interface Name: %v, IP Address: %v\n", i.Name, addr.String())
		}
	}

	// Copy configuration for the network
	genesisURL, configURL, urlSet := nu.GetGenesisAndConfigurationFromEnv()

	log.Printf("Network: %s", network)
	if !urlSet {
		log.Printf("Profile: %s", profile)
	}
	log.Printf("Overwrite Config: %t", overwriteConfig)

	if urlSet {
		log.Printf("Using genesis and configuration URLs from environment variables: %s, %s", genesisURL, configURL)
		err := nu.DownloadNetworkConfiguration(genesisURL, configURL, algodDataDir)
		if err != nil {
			fmt.Printf("Failed to download network configuration: %v", err)
			os.Exit(1)
		}
	} else {
		fu := utils.FileUtils{}
		err := fu.CopyNetworkConfigurationFromFilesystem(network, profile, overwriteConfig, genesisJSONPathFmt, configJSONPathFmt, algodDataDir)
		if err != nil {
			fmt.Printf("Failed to copy network configuration: %v", err)
			os.Exit(1)
		}
	}

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
