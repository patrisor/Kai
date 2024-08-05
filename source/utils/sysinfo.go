package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Retrieves system information and returns it as a string.
func GetSystemInfo() string {
	systemInfo := map[string]string{
		"OS":               runtime.GOOS,
		"Architecture":     runtime.GOARCH,
		"CPU Count":        fmt.Sprintf("%d", runtime.NumCPU()),
		"Go Version":       runtime.Version(),
		"Hostname":         GetHostname(),
		"Current User":     GetCurrentUser(),
		"Home Directory":   GetHomeDirectory(),
		"Environment Vars": GetEnvironmentVariables(),
	}
	info, _ := json.MarshalIndent(systemInfo, "", "  ")
	return string(info)
}

// GetHostname returns the hostname of the system.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// GetCurrentUser returns the current user's name.
func GetCurrentUser() string {
	user := os.Getenv("USER")
	if user == "" {
		return "unknown"
	}
	return user
}

// GetHomeDirectory returns the home directory of the current user.
func GetHomeDirectory() string {
	home := os.Getenv("HOME")
	if home == "" {
		return "unknown"
	}
	return home
}

// GetEnvironmentVariables returns all environment variables as a string.
func GetEnvironmentVariables() string {
	envVars := os.Environ()
	return strings.Join(envVars, "\n")
}