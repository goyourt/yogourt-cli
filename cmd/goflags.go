package cmd

import (
	"os"
	"runtime"
	"strings"
)

const stripDwarfLinkerFlag = "-ldflags=-w"

func disablePluginDwarfOnDarwin() error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return os.Setenv("GOFLAGS", withStripDwarfLinkerFlag(os.Getenv("GOFLAGS")))
}

func withStripDwarfLinkerFlag(goflags string) string {
	if strings.Contains(goflags, "-ldflags") {
		return goflags
	}
	return strings.TrimSpace(goflags + " " + stripDwarfLinkerFlag)
}
