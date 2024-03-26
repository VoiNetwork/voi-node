package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

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

		fmt.Println("Catching up to catchpoint:", catchpoint)
		cmd := exec.Command("/node/bin/goal", "-d", "/algod/data", "node", "catchup", catchpoint)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}
	} else {
		fmt.Println("last-catchpoint not found in the response")
	}
}
