package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func copyConfiguration(srcFile string, destFile string) error {
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		src, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer src.Close()

		destDir := filepath.Dir(destFile)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err = os.MkdirAll(destDir, 0755)
			if err != nil {
				return err
			}
		}

		dest, err := os.Create(destFile)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		return err
	}
	return nil
}

func executeCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute %s: %v", command, err)
	}
}

func startProcess(command string, args ...string) chan bool {
	done := make(chan bool)

	go func() {
		executeCommand(command, args...)
		done <- true
	}()

	return done
}

func main() {
	err := copyConfiguration("/algod/configuration/genesis.json", "/algod/data/genesis.json")
	if err != nil {
		log.Fatalf("Failed to copy genesis.json: %v", err)
	}

	err = copyConfiguration("/algod/configuration/config.json", "/algod/data/config.json")
	if err != nil {
		log.Fatalf("Failed to copy config.json: %v", err)
	}

	// Start algod
	done := startProcess("/node/bin/algod", "-d", "/algod/data")

	// Wait for algod to start
	time.Sleep(5 * time.Second)

	// Execute catch-catchpoint
	executeCommand("/node/bin/catch-catchpoint")

	<-done
}
