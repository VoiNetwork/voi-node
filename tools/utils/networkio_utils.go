package utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type NetworkIOUtils struct{}

func (nu NetworkIOUtils) IsPortOpen(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (nu NetworkIOUtils) LookupSRVRecords(network Network) ([]string, error) {
	_, srvs, err := net.LookupSRV("archive", "tcp", network.ArchivalDNS)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup SRV records: %v", err)
	}

	if len(srvs) == 0 {
		return nil, fmt.Errorf("no SRV records found")
	}

	var records []string
	for _, srv := range srvs {
		records = append(records, fmt.Sprintf("%s:%d", srv.Target, srv.Port))
	}

	return records, nil
}

func (nu NetworkIOUtils) DownloadNetworkConfiguration(genesisURL, algodDataDir string) error {
	if err := downloadFile(genesisURL, filepath.Join(algodDataDir, "genesis.json")); err != nil {
		return fmt.Errorf("failed to download genesis.json: %w", err)
	}

	return nil
}

func downloadFile(url, destFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	fu := FileUtils{}
	err = fu.EnsureDirExists(destFile)
	if err != nil {
		return err
	}

	out, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", destFile, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing data to %s: %w", destFile, err)
	}

	log.Printf("Successfully downloaded %s to %s", url, destFile)
	return nil
}
