/*
Copyright © 2026

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Whitelist of allowed environment variable prefixes for security
var allowedEnvPrefixes = []string{
	"mqtt_", "phev_", "log_", "ha_", "vehicle_",
	"update_", "wifi_", "remote_", "local_",
	"route_", "connection_", "availability_",
	"encoding_", "config_",
}

// isAllowedEnvVar checks if an environment variable is in the allowed list
func isAllowedEnvVar(key string) bool {
	key = strings.ToLower(key)
	for _, prefix := range allowedEnvPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// sanitizeEnvValue removes potentially dangerous characters from env values
func sanitizeEnvValue(value string) string {
	// Remove null bytes
	value = strings.ReplaceAll(value, "\x00", "")
	
	// Limit length to prevent memory exhaustion
	const maxEnvValueLength = 4096
	if len(value) > maxEnvValueLength {
		log.Warnf("SECURITY: Environment value truncated (exceeded %d bytes)", maxEnvValueLength)
		value = value[:maxEnvValueLength]
	}
	
	return value
}

// ConfigReloader handles hot reloading of configuration
type ConfigReloader struct {
	configFile     string
	lastModTime    time.Time
	checkInterval  time.Duration
	mu             sync.RWMutex
	stopChan       chan struct{}
	reloadCallback func()
}

// NewConfigReloader creates a new configuration reloader
func NewConfigReloader(configFile string, checkInterval time.Duration) *ConfigReloader {
	return &ConfigReloader{
		configFile:    configFile,
		checkInterval: checkInterval,
		stopChan:      make(chan struct{}),
	}
}

// SetReloadCallback sets the function to call when configuration is reloaded
func (r *ConfigReloader) SetReloadCallback(callback func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reloadCallback = callback
}

// Start begins watching the configuration file for changes
func (r *ConfigReloader) Start() {
	// Get initial modification time
	if info, err := os.Stat(r.configFile); err == nil {
		r.lastModTime = info.ModTime()
	}

	go r.watchLoop()
	log.Infof("Configuration hot reload enabled - watching %s", r.configFile)
}

// Stop stops watching the configuration file
func (r *ConfigReloader) Stop() {
	close(r.stopChan)
}

// watchLoop continuously checks for file changes
func (r *ConfigReloader) watchLoop() {
	ticker := time.NewTicker(r.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.checkAndReload()
		case <-r.stopChan:
			return
		}
	}
}

// checkAndReload checks if the file has changed and reloads if necessary
func (r *ConfigReloader) checkAndReload() {
	info, err := os.Stat(r.configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Debugf("Error checking config file: %v", err)
		}
		return
	}

	modTime := info.ModTime()
	
	// Only reload if the modification time has actually changed
	if modTime.Equal(r.lastModTime) || modTime.Before(r.lastModTime) {
		return
	}
	
	log.Infof("Configuration file changed (timestamp: %v), reloading...", modTime)
	r.lastModTime = modTime

	// Re-read environment variables from the file
	if err := r.reloadEnvFile(); err != nil {
		log.Errorf("Failed to reload configuration: %v", err)
		return
	}

	// Call the reload callback if set
	r.mu.RLock()
	callback := r.reloadCallback
	r.mu.RUnlock()

	if callback != nil {
		callback()
	}

	log.Infof("Configuration reloaded successfully")
}

// reloadEnvFile reads the .env file and updates environment variables
func (r *ConfigReloader) reloadEnvFile() error {
	// Read the .env file
	data, err := os.ReadFile(r.configFile)
	if err != nil {
		return err
	}

	// Parse and set environment variables
	lines := parseEnvFile(string(data))
	validCount := 0
	blockedCount := 0
	
	for key, value := range lines {
		// SECURITY: Only allow whitelisted environment variables
		if !isAllowedEnvVar(key) {
			log.Warnf("SECURITY: Ignoring unauthorized environment variable: %s", key)
			blockedCount++
			continue
		}
		
		// SECURITY: Sanitize value
		value = sanitizeEnvValue(value)
		
		// Update environment variable
		os.Setenv(key, value)
		log.Debugf("Updated %s", key) // Don't log value for security
		validCount++
	}
	
	if blockedCount > 0 {
		log.Warnf("SECURITY: Blocked %d unauthorized environment variables", blockedCount)
	}
	log.Debugf("Loaded %d valid environment variables", validCount)

	// Tell viper to re-read environment variables
	viper.AutomaticEnv()

	return nil
}

// parseEnvFile parses a .env file format
func parseEnvFile(content string) map[string]string {
	result := make(map[string]string)
	
	for _, line := range splitLines(content) {
		line = trimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}

		// Find the = sign
		eqIdx := -1
		for i := 0; i < len(line); i++ {
			if line[i] == '=' {
				eqIdx = i
				break
			}
		}

		if eqIdx == -1 {
			continue
		}

		key := trimSpace(line[:eqIdx])
		value := trimSpace(line[eqIdx+1:])

		// Remove quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		result[key] = value
	}

	return result
}

// Helper functions to avoid external dependencies
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading space
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	
	// Trim trailing space
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}

// GetConfigFilePath returns the path to the .env file
func GetConfigFilePath() string {
	// Check if running in Docker with /config mount
	dockerPath := "/config/.env"
	if _, err := os.Stat(dockerPath); err == nil {
		return dockerPath
	}

	// Check current directory
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	}

	// Check home directory
	if home, err := os.UserHomeDir(); err == nil {
		homeEnv := filepath.Join(home, ".phev2mqtt.env")
		if _, err := os.Stat(homeEnv); err == nil {
			return homeEnv
		}
	}

	// Default to Docker path
	return dockerPath
}
