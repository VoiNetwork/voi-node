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

	resp, err := client.Get("https://testnet-api.voi.nodly.io/v2/status")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
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
	} else {
		fmt.Println("HTTP request failed with status:", resp.StatusCode)
	}
}
