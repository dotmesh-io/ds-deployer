package version

import (
	"runtime"
)

const (
	ProductName string = "dotscience-deployer"
	APIVersion  string = "1"
)

// Revision that was compiled. This will be filled in by the compiler.
var Revision string

// BuildDate is when the binary was compiled.  This will be filled in by the
// compiler.
var BuildDate string

// Version number that is being run at the moment.  Version should use semver.
var Version = "dev"

// Experimental is intended to be used to enable alpha features.
var Experimental string

// VersionInfo describes version and runtime info.
type VersionInfo struct {
	Name          string `json:"name"`
	BuildDate     string `json:"buildDate"`
	Revision      string `json:"revision"`
	Version       string `json:"version"`
	APIVersion    string `json:"apiVersion"`
	GoVersion     string `json:"goVersion"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	KernelVersion string `json:"kernelVersion"`
	Experimental  bool   `json:"experimental"`
}

// GetVersion returns version info.
func GetVersion() VersionInfo {
	v := VersionInfo{
		Name:       ProductName,
		Revision:   Revision,
		BuildDate:  BuildDate,
		Version:    Version,
		APIVersion: APIVersion,
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
	return v
}
