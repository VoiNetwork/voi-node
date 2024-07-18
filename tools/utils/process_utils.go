package utils

import (
	"bytes"
	"os"
	"os/exec"
)

type ProcessUtils struct{}

func (pu ProcessUtils) ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return stdoutBuf.String(), err
}

func (pu ProcessUtils) StartProcess(command string, args ...string) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			errChan <- err
			return
		}
		close(errChan)
	}()

	return errChan
}
