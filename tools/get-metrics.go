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

const NodeExporterListenAddress = "http://localhost:9100/metrics"

func main() {
	dataDir := flag.String("d", "", "ALGORAND_DATA directory")
	flag.Parse()

	if *dataDir == "" {
		fmt.Println("Error: -d parameter is required and should point to ALGORAND_DATA")
		os.Exit(1)
	}

	client := http.Client{
		Timeout: time.Second * 5,
	}

	for {
		resp, err := client.Get(NodeExporterListenAddress)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			filePath := filepath.Join(*dataDir, "metrics.log")
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
				fmt.Printf("Successfully written %d bytes to file. HTTP status code: %d\n", bytesWritten, resp.StatusCode)
			}

			resp.Body.Close()
		}

		// Wait for 10 seconds before the next request
		time.Sleep(10 * time.Second)
	}
}
