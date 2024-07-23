package utils

import (
	"fmt"
	"io"
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

	_, err = io.Copy(file, data)
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
