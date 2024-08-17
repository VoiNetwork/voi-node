package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileUtils struct{}

func (fu FileUtils) EnsureDirExists(filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (fu FileUtils) CopyFile(srcFile, destFile string, forceOverwrite bool) error {
	if !forceOverwrite && fileExists(destFile) {
		return nil
	}

	err := fu.EnsureDirExists(destFile)
	if err != nil {
		return err
	}

	return copyFileContents(srcFile, destFile)
}

func (fu FileUtils) WriteToFile(filePath string, data io.Reader) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.Copy(file, data)
	log.Printf("Successfully written %d bytes to %s", bytes, filePath)
	if err != nil {
		return err
	}

	return nil
}

func (fu FileUtils) CopyAlgodConfigurationFromFilesystem(network string, profile string, overWriteConfig bool, configJSONPathFmt string, algodDataDir string) error {
	nu := NetworkUtils{}

	if !nu.CheckIfPredefinedNetwork(network) {
		network = "testnet"
	}

	configPath := fmt.Sprintf(configJSONPathFmt, network, profile)
	return fu.CopyFile(configPath, filepath.Join(algodDataDir, "config.json"), overWriteConfig)
}

func (fu FileUtils) CopyGenesisConfigurationFromFilesystem(network string, overWriteConfig bool, genesisJSONPathFmt string, algodDataDir string) error {
	genesisPath := fmt.Sprintf(genesisJSONPathFmt, network)
	return fu.CopyFile(genesisPath, filepath.Join(algodDataDir, "genesis.json"), overWriteConfig)
}

func (fu FileUtils) UpdateJSONAttribute(filePath, attributeName string, newValue interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Parse the JSON content
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON content: %v", err)
	}

	// Update the specified attribute with the new value
	data[attributeName] = newValue

	// Write the updated JSON content back to the file
	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated JSON content: %v", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated JSON content to file: %v", err)
	}

	return nil
}

func (fu FileUtils) EnsureGUIDExists(filePath string) error {
	// Open and read the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Parse the JSON content
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON content: %v", err)
	}

	// Check if the GUID attribute is empty and generate a new GUID if necessary
	if guid, ok := data["GUID"].(string); ok && guid == "" {
		newGUID := uuid.New().String()
		data["GUID"] = newGUID
	}

	// Write the updated JSON content back to the file
	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated JSON content: %v", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated JSON content to file: %v", err)
	}

	return nil
}

func (fu FileUtils) SetTelemetryState(filePath, telemetryName string, enabled bool) error {
	err := fu.UpdateJSONAttribute(filePath, "Enable", enabled)
	if err != nil {
		return fmt.Errorf("failed to set telemetry enabled state: %v", err)
	}
	err = fu.UpdateJSONAttribute(filePath, "Name", telemetryName)
	if err != nil {
		return fmt.Errorf("failed to set telemetry name: %v", err)
	}

	err = fu.EnsureGUIDExists(filePath)
	if err != nil {
		return fmt.Errorf("failed to ensure GUID exists: %v", err)
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func copyFileContents(srcFile string, destFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}
