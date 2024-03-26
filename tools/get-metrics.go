package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const NodeExporterListenAddress = "http://relay:9100/metrics"

func retrieveMetrics(dataDir *string) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	for {
		resp, err := client.Get(NodeExporterListenAddress)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			filePath := filepath.Join(*dataDir, "algod.prom")
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			bytesWritten, err := io.Copy(file, resp.Body)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			} else {
				fmt.Printf("Successfully written %d bytes to %s. HTTP status code: %d\n", bytesWritten, filePath, resp.StatusCode)
			}

			resp.Body.Close()
		}

		// Wait for 10 seconds before the next request
		time.Sleep(10 * time.Second)
	}
}

func main() {
	dataDir := flag.String("d", "", "ALGORAND_DATA directory")
	flag.Parse()

	if *dataDir == "" {
		fmt.Println("Error: -d parameter is required and should point to ALGORAND_DATA")
		os.Exit(1)
	}

	err := os.MkdirAll(*dataDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	retrieveMetrics(dataDir)
}
