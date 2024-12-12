// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package main

// package main is an executable for opentsg job control in minikube

import (
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/cli"
)

// release history is in the cli package
// main sets the logger and runs the command line interface (cli)
func main() {

	// run the cli processor
	cli.Execute()
}
