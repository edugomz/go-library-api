// Package version holds the build-time version string.
package version

// Version is overridden at build time via:
//
//	go build -ldflags "-X library-api/internal/version.Version=$(git describe --tags --always --dirty)"
var Version = "dev"
