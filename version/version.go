package version

import (
	"fmt"
	"runtime"
)

var (
	// application's name
	Name = ""
	// application's version string
	Version = ""
	// commit
	Commit = ""
)

// Info defines the application version information.
type Info struct {
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"commit" yaml:"commit"`
	GoVersion string `json:"go" yaml:"go"`
}

func (v Info) String() string {
	return fmt.Sprintf(`%s: %s
git commit: %s
%s`, v.Name, v.Version, v.GitCommit, v.GoVersion)
}

func NewInfo() Info {
	return Info{
		Name:      Name,
		Version:   Version,
		GitCommit: Commit,
		GoVersion: fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)}
}
