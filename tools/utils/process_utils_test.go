package utils

import (
	"testing"
	"time"
)

func TestExecuteCommand(t *testing.T) {
	pu := ProcessUtils{}

	t.Run("Success", func(t *testing.T) {
		_, err := pu.ExecuteCommand("echo", "Hello, world!")
		if err != nil {
			t.Errorf("ExecuteCommand() with echo should not have an error, got %v", err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		_, err := pu.ExecuteCommand("command-that-does-not-exist")
		if err == nil {
			t.Error("ExecuteCommand() with a non-existent command should have an error, got nil")
		}
	})
}

func TestStartProcess(t *testing.T) {
	pu := ProcessUtils{}

	t.Run("Success", func(t *testing.T) {
		errChan := pu.StartProcess("echo", "Hello, world!")
		select {
		case err, ok := <-errChan:
			if !ok {
				// Channel closed without errors, which is expected
			} else if err != nil {
				t.Errorf("StartProcess() with echo should not have an error, got %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("StartProcess() timed out, expected completion")
		}
	})

	t.Run("Fail", func(t *testing.T) {
		errChan := pu.StartProcess("command-that-does-not-exist")
		select {
		case err, ok := <-errChan:
			if !ok {
				t.Error("StartProcess() with a non-existent command should not close the channel without sending an error")
			} else if err == nil {
				t.Error("StartProcess() with a non-existent command should have an error")
			}
		case <-time.After(5 * time.Second):
			t.Error("StartProcess() timed out, expected an error")
		}
	})
}
