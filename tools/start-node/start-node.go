package main

import (
	"flag"
	"github.com/voinetwork/voi-node/tools/utils"
	"log"
	"os"
	"time"
)

const (
	envCatchupVar = "VOINETWORK_CATCHUP"
	algodDataDir  = "/algod/data"
	algodCmd      = "/node/bin/algod"
	catchupCmd    = "/node/bin/catch-catchpoint"
	goalCmd       = "/node/bin/goal"
)

var network string
var profile string
var overwriteConfig bool

func init() {
	flag.StringVar(&network, "network", "testnet", "Specify the network (testnet)")
	flag.StringVar(&profile, "profile", "relay", "Specify the profile (archiver, relay, developer)")
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

	genesisURL, urlSet := nu.GetGenesisFromEnv()

	envProfile, profileSet := nu.GetProfileFromEnv()
	if profileSet {
		profile = envProfile
	}

	log.Printf("Network: %s", network)
	log.Printf("Profile: %s", profile)
	log.Printf("Overwrite Config: %t", overwriteConfig)

	cu := utils.ConfigUtils{}
	cu.HandleConfiguration(urlSet, genesisURL, network, profile, overwriteConfig, algodDataDir)

	pu := utils.ProcessUtils{}
	var done <-chan error

	envCatchup := os.Getenv(envCatchupVar)
	if profile == "archiver" && envCatchup != "0" {
		predefinedNetwork, err := nu.NewNetwork(network)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		niou := utils.NetworkIOUtils{}
		srvRecords, err := niou.LookupSRVRecords(predefinedNetwork)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("Catching up using direct archiver connection to %s: ", srvRecords[0])
		done = pu.StartProcess(goalCmd, "node", "start", "-d", algodDataDir, "-p", srvRecords[0])
	} else {

		done = pu.StartProcess(algodCmd, "-d", algodDataDir)

		if envCatchup != "0" && !urlSet && profile != "archiver" {
			retryCount := 0
			maxRetries := 10

			// allow algod to start up before executing catchup command
			time.Sleep(5 * time.Second)
			for retryCount < maxRetries {
				_, err := pu.ExecuteCommand(catchupCmd)
				if err == nil {
					break
				}
				retryCount++
				log.Printf("Retry %d/%d: Failed to execute catchup command, retrying in 5 seconds...", retryCount, maxRetries)
				time.Sleep(5 * time.Second)
			}
			if retryCount == maxRetries {
				log.Printf("Failed to execute catchup command after %d retries", maxRetries)
				return
			}
		} else {
			log.Printf("Skipping catchup execution as %s is set to 0", envCatchupVar)
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Fatalf("Process finished with error: %v", err)
			}
		case <-ticker.C:
			// do nothing
		}
	}
}
