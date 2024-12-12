//  Copyright ©2019-2024    Mr MXF   info@mrmxf.com
//  BSD-3-Clause License    https://opensource.org/license/bsd-3-clause/
//
// Package cli provides a simple command line interface to launch the
// watchfolder node controller. It is based on package github.com/spf13/cobra.

package cli

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/job"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/log"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/semver"
)

//go:embed releases.yaml
var eFs embed.FS

var jobsFolder = "."
var LogLevelDebug = false
var ProductionLogging = false

// mainCmd starts the controller after flags have been parsed.
var mainCmd = &cobra.Command{
	Use:   "opentsg-ctl-watchfolder",
	Short: "run a watchfolder controller in the current folder",
	Long:  `watch a folder for opentsg subfolder jobs  named "jobNNNN" where N is [0-9].`,

	// Run starts the watchfolder controller. It never stop unless killed.
	Run: func(root *cobra.Command, args []string) {

		// if we have to show the version, do it and exit
		if ShowVersion || ShowVersionShort || ShowVersionNote {
			slog.Debug("root command: show version command")
			versionCmd.Run(root, args)
			return
		}

		// set the log level based on flags
		switch {
		case LogLevelDebug && ProductionLogging:
			log.UseProductionJSONErrorLogger()
		case (!LogLevelDebug) && ProductionLogging:
			log.UseJSONInfoLogger()
		case LogLevelDebug && (!ProductionLogging):
			log.UsePrettyDebugLogger()
		case (!LogLevelDebug) && (!ProductionLogging):
			log.UsePrettyInfoLogger()
		}

		startMsg := fmt.Sprintf("Minikube Watchfolder Controller (%s)", jobsFolder)
		slog.Info(startMsg)

		//initialise asn empty list of jobs
		jobs := &job.JobManagement{
			Folder:       jobsFolder,
			LockFileName: "_status.lock",
			JobLogName:   "_ctl-watchfolder.log",
		}

		//init the jobs to fast start the polling
		jobs.ParseJobs()
		jobs.HandleJobs()
		//polling loop
		jobs.Wg.Add(1)
		jobs.StartPolling()
		// if the polling goroutine crashes the end the executable and wait for
		// minikube to restart it
		jobs.Wg.Wait()
	},
}

func Execute() {
	if err := mainCmd.Execute(); err != nil {
		slog.Error("Failed to initialise command line interface", "err", err)
		os.Exit(1)
	}
}

func init() {
	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	//initilise the version history
	err := semver.Initialise(eFs, "releases.yaml")
	if err != nil {
		slog.Debug("init semver failed", "err", err)
	}

	// jobsFolder flags
	mainCmd.PersistentFlags().StringVarP(&jobsFolder, "folder", "f", ".", "watch a folder for new jobs e.g. --folder \"./network-jobs/\"")

	// initialise the version flags
	mainCmd.PersistentFlags().BoolVar(&ShowVersion, "version", false, "show full semantic version, program name etc.")
	mainCmd.PersistentFlags().BoolVarP(&ShowVersionShort, "v", "v", false, "show short semantic version")
	mainCmd.PersistentFlags().BoolVarP(&ShowVersionNote, "note", "n", false, "show just the version note")

	// logging flags
	mainCmd.PersistentFlags().BoolVarP(&LogLevelDebug, "debug", "D", false, "set logging level to debug (or info production mode)")
	mainCmd.PersistentFlags().BoolVarP(&LogLevelDebug, "production", "P", false, "production mode - JSON logging at error / info level")

	// config file flags
	// mainCmd.PersistentFlags().StringVarP(&clCmd.ConfigFilePath, "config", "c", "", "clog -c myClogfig.yaml   # clog Core Cat clogrc/core/clog.clConfig.yaml > myClogfig.yaml")
}
