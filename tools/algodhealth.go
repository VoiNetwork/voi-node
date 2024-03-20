package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get("http://127.0.0.1:8080/health")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP request failed with status: %d\n", resp.StatusCode)
		os.Exit(1)
	} else {
		fmt.Printf("Algod healthcheck passed.")
	}
}
