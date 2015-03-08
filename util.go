package main

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetUserCacheDirectory() string {
	if baseDir := os.Getenv("XDG_CACHE_HOME"); baseDir != "" {
		return filepath.Join(baseDir, "rfc")
	}

	if user, err := user.Current(); err == nil {
		return filepath.Join(user.HomeDir, ".cache", "rfc")
	}

	if homeDir := os.Getenv("HOME"); homeDir != "" {
		return filepath.Join(homeDir, ".cache", "rfc")
	}

	return ""
}
