package xdgdirs

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/redforks/hal"
)

const tag = "xdgdirs"

// Home return current user Home directory
func Home() string {
	r := hal.Getenv("HOME")
	if r == "" {
		log.Panicf("[%s] HOME environment string not defined", tag)
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

// ConfigHome() is directory stores user specific configuration files.
func ConfigHome() string {
	return getDirWithDefault("XDG_CONFIG_HOME", ".config")
}

// CacheHome() is a directory that user specific non-essential (cached) data
// should be written.
func CacheHome() string {
	return getDirWithDefault("XDG_CACHE_HOME", ".cache")
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
