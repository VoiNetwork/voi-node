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
	// Define and parse the -d flag
	dataDir := flag.String("d", ".", "ALGORAND_DATA directory")
	flag.Parse()

	client := http.Client{
		Timeout: time.Second * 5,
	}

	for {
		resp, err := client.Get(NodeExporterListenAddress)
		if err != nil {
			// Log the error and continue
			fmt.Println("Error:", err)
		} else {
			// Create the metrics.log file in the provided directory
			filePath := filepath.Join(*dataDir, "metrics.log")
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			// Copy the response body to the file
			bytesWritten, err := io.Copy(file, resp.Body)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			} else {
				fmt.Printf("Successfully written %d bytes to file. HTTP status code: %d\n", bytesWritten, resp.StatusCode)
			}

			// Close the response body
			resp.Body.Close()
		}

		// Wait for 10 seconds before the next request
		time.Sleep(10 * time.Second)
	}
}
