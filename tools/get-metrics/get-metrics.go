package main

import (
	"flag"
	"github.com/voinetwork/voi-node/tools/utils"
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

var metricsDir string
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func init() {
	flag.StringVar(&metricsDir, "d", "", "Specify the metrics directory")
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

	return fu.WriteToFile(filePath, resp.Body)
}

func main() {
	flag.Parse()

	if metricsDir == "" {
		log.Println("Error: -d parameter is required and should point to metrics directory")
		os.Exit(1)
	}

	err := os.MkdirAll(metricsDir, 0755)
	if err != nil {
		log.Println("Error creating directory:", err)
		os.Exit(1)
	}

	if err := retrieveMetrics(metricsDir); err != nil {
		log.Println("Error retrieving metrics:", err)
		os.Exit(1)
	}
}
