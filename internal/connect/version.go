package connect

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	version string
)

// GetShortenedVersion returns the short program version
func GetShortenedVersion() string {
	return strings.Split(version, "~")[0]
}


// GetFullVersion returns the full program version
func GetFullVersion() string {
	return version
}