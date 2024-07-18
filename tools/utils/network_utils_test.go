package utils

import (
	"os"
	"testing"
)

func TestGetStatusURL(t *testing.T) {
	tests := []struct {
		name    string
		network string
		wantURL string
		wantErr bool
	}{
		{"testNet", "testnet", "https://testnet-api.voi.nodly.io/v2/status", false},
		{"Unknown", "unknown", "", true},
	}

	nu := NetworkUtils{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := nu.GetStatusURL(tt.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatusURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotURL != tt.wantURL {
				t.Errorf("GetStatusURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestNetworkUtils_GetEnvNetworkVar(t *testing.T) {
	expected := envNetworkVar
	nu := NetworkUtils{}
	if got := nu.GetEnvNetworkVar(); got != expected {
		t.Errorf("GetEnvNetworkVar() = %v, want %v", got, expected)
	}
}

func TestNetworkUtils_GetNetworkFromEnv(t *testing.T) {
	nu := NetworkUtils{}

	// Test with environment variable set
	expectedNetwork := "testnet"
	os.Setenv(envNetworkVar, expectedNetwork)
	defer os.Unsetenv(envNetworkVar)

	got, ok := nu.GetNetworkFromEnv()
	if !ok || got != expectedNetwork {
		t.Errorf("GetNetworkFromEnv() = %v, %v, want %v, true", got, ok, expectedNetwork)
	}

	// Test without environment variable set
	os.Unsetenv(envNetworkVar)
	_, ok = nu.GetNetworkFromEnv()
	if ok {
		t.Error("GetNetworkFromEnv() expected to return false when environment variable is not set")
	}
}
