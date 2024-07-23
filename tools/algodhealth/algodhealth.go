package main

import (
	"github.com/voinetwork/voi-node/tools/utils"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	algodDataDir   = "/algod/data"
	goalCmd        = "/node/bin/goal"
	healthCheckURL = "http://localhost:8080/health"
	httpTimeout    = 5 * time.Second
)

func checkAlgodProcess() bool {
	pu := utils.ProcessUtils{}
	_, err := pu.ExecuteCommand(goalCmd, "node", "status", "-d", algodDataDir)
	return err == nil
}

func checkAlgodHTTPHealth() bool {
	client := createHTTPClient()
	resp, err := client.Get(healthCheckURL)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed with status: %d", resp.StatusCode)
		return false
	}
	return true
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

func main() {
	if checkAlgodProcess() || checkAlgodHTTPHealth() {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("Failed to get hostname: %v", err)
			return
		}
		log.Printf("Algod healthcheck passed on host: %s\n", hostname)
	} else {
		os.Exit(1)
	}
}
