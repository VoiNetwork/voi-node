package main

import (
	"encoding/json"
	"flag"
	"github.com/voinetwork/voi-node/tools/utils"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	envNetworkVar     = "VOINETWORK_NETWORK"
	goalCmd           = "/node/bin/goal"
	algodDataDir      = "/algod/data"
	httpTimeout       = 30 * time.Second
	httpRetryAttempts = 10
)

var networkArgument string

func init() {
	flag.StringVar(&networkArgument, "network", "testnet", "Specify the network (testnet)")
}

func getLastNodeRound(pu utils.ProcessUtils) (int, error) {
	output, err := pu.ExecuteCommand(goalCmd, "node", "lastround", "-d", algodDataDir)

	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	outputStr := strings.TrimSpace(output)
	lastRound, err := strconv.Atoi(outputStr)
	if err != nil {
		log.Fatalf("Error parsing output: %v", err)
	}

	return lastRound, nil
}

func main() {
	flag.Parse()
	envNetwork := os.Getenv(envNetworkVar)
	if envNetwork != "" {
		log.Printf("Using network from environment variable: %s", envNetwork)
		networkArgument = envNetwork
	}

	nu := utils.NetworkUtils{}
	network, err := nu.NewNetwork(networkArgument)

	log.Printf("Catchup on network: %s", network.Name)

	pu := utils.ProcessUtils{}
	statusURL := network.StatusURL
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	client := http.Client{Timeout: httpTimeout}
	var resp *http.Response

	for i := 0; i < httpRetryAttempts; i++ {
		resp, err = client.Get(statusURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Printf("Attempt %d failed with error: %v. Retrying in %v...\n", i+1, err, waitTime)
		time.Sleep(waitTime)
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get a successful response: %v", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	var result map[string]interface{}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return
	}

	if lastCatchpoint, ok := result["last-catchpoint"]; ok {
		catchpoint, ok := lastCatchpoint.(string)
		if !ok {
			log.Fatal("Error: last-catchpoint is not a string")
			return
		}

		var catchpointRound int
		catchpointParts := strings.Split(catchpoint, "#")
		if len(catchpointParts) > 0 {
			catchpointRound, err = strconv.Atoi(catchpointParts[0])
			if err != nil {
				log.Fatal("Error: catchpoint round is not an integer")
				return
			}
			catchpointRoundStr := strconv.Itoa(catchpointRound)
			log.Printf("Catchpoint round: %s", catchpointRoundStr)
		} else {
			log.Fatal("Error: catchpoint does not contain '#'")
		}

		lastNodeRound, _ := getLastNodeRound(pu)

		var lastNetworkRound int
		if lastRound, ok := result["last-round"].(float64); ok {
			lastNetworkRound = int(lastRound)
		} else {
			log.Println("last-round not found in the response or is not a float64")
			return
		}

		log.Printf("Last node round: %d, Last network round: %d\n", lastNodeRound, lastNetworkRound)

		if (lastNodeRound) > lastNetworkRound-1000 {
			log.Print("Current round is not that far behind (if at all), skipping catchup")
			return
		} else if catchpointRound < lastNodeRound-1000 {
			log.Print("Catchpoint round is behind the network, skipping catchup")
			return
		}

		log.Printf("Catching up to catchpoint: %s", catchpoint)
		_, err = pu.ExecuteCommand(goalCmd, "-d", algodDataDir, "node", "catchup", catchpoint)
		if err != nil {
			log.Fatalf("Error running command: %v", err)
			return
		}
	} else {
		log.Print("last-catchpoint not found in the response")
	}
}
