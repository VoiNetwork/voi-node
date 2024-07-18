package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileUtils struct{}

func (fu FileUtils) EnsureDirExists(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (fu FileUtils) CopyFile(srcFile string, destFile string, forceOverwrite bool) error {

	if _, err := os.Stat(destFile); err == nil {
		if !forceOverwrite {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	err := fu.EnsureDirExists(filepath.Dir(destFile))
	if err != nil {
		return err
	}

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

func (fu FileUtils) WriteToFile(filePath string, data io.Reader, statusCode int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytesWritten, err := io.Copy(file, data)
	if err != nil {
		return err
	}

	log.Printf("Successfully written %d bytes to %s. HTTP status code: %d\n", bytesWritten, filePath, statusCode)
	return nil
}

func (fu FileUtils) CopyNetworkConfigurationFromFilesystem(network string, profile string, overWriteConfig bool, genesisJSONPathFmt string, configJSONPathFmt string, algodDataDir string) error {
	err := fu.CopyFile(fmt.Sprintf(genesisJSONPathFmt, network), algodDataDir+"/genesis.json", overWriteConfig)
	if err != nil {
		return fmt.Errorf("failed to copy genesis.json: %v", err)
	}

	err = fu.CopyFile(fmt.Sprintf(configJSONPathFmt, network, profile), algodDataDir+"/config.json", overWriteConfig)
	if err != nil {
		return fmt.Errorf("failed to copy config.json: %v", err)
	}
	return nil
}
