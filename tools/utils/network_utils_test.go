package utils

import (
	"os"
	"testing"
)

func TestNewNetwork(t *testing.T) {
	nu := NetworkUtils{}

	network, err := nu.NewNetwork(testNet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if network.Name != testNet {
		t.Errorf("Expected network name %s, got %s", testNet, network.Name)
	}

	_, err = nu.NewNetwork("invalid")
	if err == nil {
		t.Fatalf("Expected error for unsupported network, got none")
	}
}

func TestCheckIfPredefinedNetwork(t *testing.T) {
	nu := NetworkUtils{}

	if !nu.CheckIfPredefinedNetwork(testNet) {
		t.Errorf("Expected true for predefined network %s", testNet)
	}

	if nu.CheckIfPredefinedNetwork("invalid") {
		t.Errorf("Expected false for unsupported network")
	}
}

func TestGetEnvNetworkVar(t *testing.T) {
	nu := NetworkUtils{}
	expected := envNetworkVar

	if nu.GetEnvNetworkVar() != expected {
		t.Errorf("Expected %s, got %s", expected, nu.GetEnvNetworkVar())
	}
}

func TestGetNetworkFromEnv(t *testing.T) {
	nu := NetworkUtils{}
	expected := "test_network"
	os.Setenv(envNetworkVar, expected)
	defer os.Unsetenv(envNetworkVar)

	network, found := nu.GetNetworkFromEnv()
	if !found {
		t.Fatalf("Expected to find network from environment variable")
	}
	if network != expected {
		t.Errorf("Expected %s, got %s", expected, network)
	}
}

func TestGetEnvProfileVar(t *testing.T) {
	nu := NetworkUtils{}
	expected := envProfileVar

	if nu.GetEnvProfileVar() != expected {
		t.Errorf("Expected %s, got %s", expected, nu.GetEnvProfileVar())
	}
}

func TestGetProfileFromEnv(t *testing.T) {
	nu := NetworkUtils{}
	expected := "test_profile"
	os.Setenv(envProfileVar, expected)
	defer os.Unsetenv(envProfileVar)

	profile, found := nu.GetProfileFromEnv()
	if !found {
		t.Fatalf("Expected to find profile from environment variable")
	}
	if profile != expected {
		t.Errorf("Expected %s, got %s", expected, profile)
	}
}

func TestGetGenesisFromEnv(t *testing.T) {
	nu := NetworkUtils{}
	expected := "test_genesis_url"
	os.Setenv(envGenesisURLVar, expected)
	defer os.Unsetenv(envGenesisURLVar)

	genesisURL, found := nu.GetGenesisFromEnv()
	if !found {
		t.Fatalf("Expected to find genesis URL from environment variable")
	}
	if genesisURL != expected {
		t.Errorf("Expected %s, got %s", expected, genesisURL)
	}
}
