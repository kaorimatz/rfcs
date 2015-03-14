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

func GetUserCacheDirectory(appName string) string {
	if baseDir := os.Getenv("XDG_CACHE_HOME"); baseDir != "" {
		return filepath.Join(baseDir, appName)
	}

	if user, err := user.Current(); err == nil {
		return filepath.Join(user.HomeDir, ".cache", appName)
	}

	if homeDir := GetHomeDirectory(); homeDir != "" {
		return filepath.Join(homeDir, ".cache", appName)
	}

	return ""
}
