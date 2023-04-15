package internal

import (
	"fmt"
)

var (
	version    = "dev"
	commitHash = "n/a"
	buildTime  = "n/a"
)

type Version struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildTime  string `json:"build_time"`
}

var versionData = Version{
	Version:    version,
	CommitHash: commitHash,
	BuildTime:  buildTime,
}

func BuildVersionString() string {
	return fmt.Sprintf(
		"Version: %s\nCommit hash: %s\nBuild time: %s",
		versionData.Version,
		versionData.CommitHash,
		versionData.BuildTime,
	)
}

func BuildVersion() Version {
	return versionData
}
