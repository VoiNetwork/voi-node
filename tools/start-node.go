package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

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

func copyConfiguration(srcFile string, destFile string) error {
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

func main() {
	err := copyConfiguration("/algod/configuration/genesis.json", "/algod/data/genesis.json")
	if err != nil {
		log.Fatalf("Failed to copy genesis.json: %v", err)
	}

	err = copyConfiguration("/algod/configuration/config.json", "/algod/data/config.json")
	if err != nil {
		log.Fatalf("Failed to copy config.json: %v", err)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(err)
			continue
		}

		for _, addr := range addrs {
			fmt.Printf("Interface Name: %v, IP Address: %v\n", i.Name, addr.String())
		}
	}

	// Start algod
	done := startProcess("/node/bin/algod", "-d", "/algod/data")

	// Wait for algod to start
	time.Sleep(5 * time.Second)

	// Execute catch-catchpoint
	executeCommand("/node/bin/catch-catchpoint")

	<-done
}
