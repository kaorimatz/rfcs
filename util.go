package main

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func GetHomeDirectory() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func GetUserCacheDirectory() string {
	if baseDir := os.Getenv("XDG_CACHE_HOME"); baseDir != "" {
		return filepath.Join(baseDir, "rfc")
	}

	if user, err := user.Current(); err == nil {
		return filepath.Join(user.HomeDir, ".cache", "rfc")
	}

	if homeDir := GetHomeDirectory(); homeDir != "" {
		return filepath.Join(homeDir, ".cache", "rfc")
	}

	return ""
}
