package xdgdirs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/redforks/errors"
	"github.com/redforks/hal"
)

const tag = "xdgdirs"

// Home return current user Home directory
func Home() string {
	r := hal.Getenv("HOME")
	if r == "" {
		r = "/"
	}
	return r
}

func getDirWithDefault(envName, def string) string {
	r := hal.Getenv(envName)
	if r != "" {
		return r
	}

	return filepath.Join(Home(), def)
}

// DataHome returns the directory that user data should written to.
func DataHome() string {
	return getDirWithDefault("XDG_DATA_HOME", ".local/share")
}

// ConfigHome is directory stores user specific configuration files.
func ConfigHome() string {
	return getDirWithDefault("XDG_CONFIG_HOME", ".config")
}

// CacheHome is a directory that user specific non-essential (cached) data
// should be written.
func CacheHome() string {
	return getDirWithDefault("XDG_CACHE_HOME", ".cache")
}

// RuntimeHome returns xdg runtime directory. Use XDG_RUNTIME_DIR if defined.
// Default to /tmp/[user]/ for non-root users, /run for root user.
//
// NOTE: /tmp/[user]/ noramlly not exist, RuntimeHome() do not create or check
// it exists.
func RuntimeHome() string {
	r := hal.Getenv("XDG_RUNTIME_DIR")
	if r != "" {
		return r
	}

	u, err := hal.CurrentUser()
	if err != nil {
		log.Panic(err)
	}

	if u.Uid == "0" { // root
		return "/run"
	}

	return filepath.Join(os.TempDir(), u.Name)
}

func getDirs(envName, homeDir, defDirs string) []string {
	env := hal.Getenv(envName)
	if env == "" {
		env = defDirs
	}

	return append([]string{homeDir}, strings.Split(env, ":")...)
}

// DataDirs returns a list of directories that data files should be find in this
// order. DataHome() always be the first entry.
func DataDirs() []string {
	return getDirs("XDG_DATA_DIRS", DataHome(), "/usr/local/share:/usr/share")
}

// ConfigDirs returns a list of directories that config files should be find in
// this order. ConfigHome() always be the first entry.
func ConfigDirs() []string {
	return getDirs("XDG_CONFIG_DIRS", ConfigHome(), "/etc/xdg")
}

// return non nil error if the path exist, but not a regular file.
func existAndRegularFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, nil
	}

	if !fileInfo.Mode().IsRegular() {
		return false, errors.Runtimef("[%s] \"%s\" exist but not regular file", tag, path)
	}

	return true, nil
}

func resolveFileInDirs(path string, dirs []string, errMsg string) (string, error) {
	for _, dir := range dirs {
		p := filepath.Join(dir, path)
		found, err := existAndRegularFile(p)
		if found {
			return p, nil
		}

		if err != nil {
			return "", err
		}
	}

	return "", errors.Runtimef(errMsg, tag, path)
}

// ResolveDataFile find file in DataDirs(), returns first found file.
// Return error if path not found, or the found path is not a regular file.
func ResolveDataFile(path string) (string, error) {
	return resolveFileInDirs(path, DataDirs(), "[%s] Can not found data file: %s")
}

// ResolveConfigFile find file in ConfigDirs(), returns first found file.
// Return error if path not found, or the found path is not a regular file.
func ResolveConfigFile(path string) (string, error) {
	return resolveFileInDirs(path, ConfigDirs(), "[%s] Can not found config file: %s")
}
