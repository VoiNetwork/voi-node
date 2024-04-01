package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("/node/bin/goal", "node", "status", "-d", "/algod/data")
	err1 := cmd.Run()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err1 == nil {
		fmt.Printf("Algod healthcheck passed on host: %s", hostname)
		os.Exit(0)
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err2 := client.Get("http://localhost:8080/health")
	if err2 != nil {
		fmt.Printf("Goal error: %v, HTTP error: %v\n", err1, err2)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP request failed with status: %d\n", resp.StatusCode)
		os.Exit(1)
	} else {
		fmt.Printf("Algod healthcheck passed on host: %s", hostname)
	}
}
