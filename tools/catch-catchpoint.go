package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getLastNodeRound() (int, error) {
	cmd := exec.Command("/node/bin/goal", "node", "lastround", "-d", "/algod/data")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return -1, err
	}

	outputStr := strings.TrimSpace(string(output))
	lastRound, err := strconv.Atoi(outputStr)
	if err != nil {
		fmt.Println("Error parsing output:", err)
		return -1, err
	}

	return lastRound, nil
}

func main() {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	var resp *http.Response
	var err error

	for i := 0; i < 10; i++ {
		resp, err = client.Get("https://testnet-api.voi.nodly.io/v2/status")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		fmt.Printf("Attempt %d failed with error: %v. Retrying in 10 seconds...\n", i+1, err)
		time.Sleep(10 * time.Second)
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", err)
		if resp != nil {
			fmt.Println("HTTP request failed with status:", resp.StatusCode)
		}
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(bodyBytes, &result)

	if lastCatchpoint, ok := result["last-catchpoint"]; ok {
		catchpoint, ok := lastCatchpoint.(string)
		if !ok {
			fmt.Println("Error: last-catchpoint is not a string")
			return
		}

		var catchpointRound int
		catchpointParts := strings.Split(catchpoint, "#")
		if len(catchpointParts) > 0 {
			catchpointRound, err = strconv.Atoi(catchpointParts[0])
			if err != nil {
				fmt.Println("Error: catchpoint round is not an integer")
				return
			}
			catchpointRoundStr := strconv.Itoa(catchpointRound)
			fmt.Println("Catchpoint round:", catchpointRoundStr)
		} else {
			fmt.Println("Error: catchpoint does not contain '#'")
		}

		lastNodeRound, err := getLastNodeRound()
		if err != nil {
			fmt.Println("Error getting last round:", err)
			return
		}

		var lastNetworkRound int

		if lastRound, ok := result["last-round"].(float64); ok {
			lastNetworkRound = int(lastRound)
		} else {
			fmt.Println("last-round not found in the response or is not a float64")
			return
		}

		fmt.Printf("Last node round: %d, Last network round: %d\n", lastNodeRound, lastNetworkRound)

		if (lastNodeRound) > lastNetworkRound-1000 {
			fmt.Println("Current round is not that far behind (if at all), skipping catchup")
			return
		} else if catchpointRound < lastNodeRound-1000 {
			fmt.Println("Catchpoint round is behind the network, skipping catchup")
			return
		}

		fmt.Println("Catching up to catchpoint:", catchpoint)
		cmd := exec.Command("/node/bin/goal", "-d", "/algod/data", "node", "catchup", catchpoint)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}
	} else {
		fmt.Println("last-catchpoint not found in the response")
	}
}
