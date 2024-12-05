//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package cmd implements commands for the cobra CLI library

package shell

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// GitCmd represents the base core
func GetShellPath() string {
	// figure out what shell we will run and log it for debugging
	whichShell, _ := exec.Command("which", "bash").CombinedOutput()
	if len(whichShell) < 3 {
		whichShell, _ = exec.Command("which", "zsh").CombinedOutput()
	}
	if len(whichShell) < 3 {
		whichShell, _ = exec.Command("which", "sh").CombinedOutput()
	}
	if len(whichShell) < 3 {
		slog.Error("Unable to find a compatible shell to run, exiting")
		os.Exit(1)
	}
	shellPath := strings.TrimSpace(string(whichShell))
	slog.Debug("Using shell: " + shellPath)

	return shellPath
}
