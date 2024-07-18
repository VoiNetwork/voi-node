package main

import (
	"flag"
	"github.com/voinetwork/docker-relay-node/tools/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	nodeExporterListenAddress = "http://relay:9100/metrics"
	httpRetryInterval         = 10 * time.Second
)

var algorandDataDir string
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func init() {
	flag.StringVar(&algorandDataDir, "d", "", "Specify the ALGORAND_DATA directory")
}

func retrieveMetrics(dataDir string) error {
	for {
		err := fetchAndStoreMetrics(dataDir)
		if err != nil {
			log.Println("Error:", err)
		}
		time.Sleep(httpRetryInterval)
	}
}

func fetchAndStoreMetrics(dataDir string) error {
	resp, err := httpClient.Get(nodeExporterListenAddress)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(dataDir, "algod.prom")
	fu := utils.FileUtils{}

	return fu.WriteToFile(filePath, resp.Body, resp.StatusCode)
}

func main() {
	flag.Parse()

	if algorandDataDir == "" {
		log.Println("Error: -d parameter is required and should point to ALGORAND_DATA")
		os.Exit(1)
	}

	err := os.MkdirAll(algorandDataDir, 0755)
	if err != nil {
		log.Println("Error creating directory:", err)
		os.Exit(1)
	}

	if err := retrieveMetrics(algorandDataDir); err != nil {
		log.Println("Error retrieving metrics:", err)
		os.Exit(1)
	}
}
