// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package main

// package main is an executable for opentsg job control in minikube

import (
	"embed"

	"github.com/opentsg/opentsg-ctl-watchfolder/cli"
)

//go:embed releases.yaml www dash/templates
var embedFs embed.FS

// main sets the logger and runs the command line interface (cli)
func main() {

	// run the cli processor
	cli.Execute(embedFs)
}
