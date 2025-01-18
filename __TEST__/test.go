// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package main

// package main creates test folders for checking logic.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/opentsg/opentsg-ctl-watchfolder/job"
	"github.com/opentsg/opentsg-ctl-watchfolder/log"
	dCopy "github.com/otiai10/copy"
	"github.com/phsym/console-slog"
)

func setupTestLogger() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)
}

type initState struct {
	id     int
	status string
	uri    string
	tpl    *string
}

var (
	tBlank = "blank"
)
var jobFolder = []initState{
	initState{id: 123, status: "NEW", tpl: nil},
	initState{id: 124, status: "QUEUED", tpl: &tBlank},
	initState{id: 125, status: "NEW test", tpl: &tBlank},
	initState{id: 126, status: "FAILED", tpl: nil},
	initState{id: 127, status: "RUNNING", tpl: nil},
	initState{id: 155, status: "QUEUED", tpl: &tBlank},
	initState{id: 157, status: "CANCELLED", tpl: nil},
	initState{id: 1538, status: "COMPLETED", tpl: nil},
}

var root = "./jobs"

func splatJobStatus(lockfile string, status string) {
	path := filepath.Join(string(lockfile), "_status.lock")
	s := []byte(status)

	if err := os.WriteFile(path, s, 0644); err != nil {
		slog.Error(fmt.Sprintf("cannot write to %s: %s", lockfile, err))
		return
	}
}

func main() {
	log.UsePrettyDebugLogger()
	slog.Info("__TEST__ reset job stats")

	jobs := &job.JobManagement{
		Folder:       root,
		LockFileName: "_status.lock",
	}

	//make all the folders & status files
	for i, f := range jobFolder {
		folderName := fmt.Sprintf("job%04d", f.id)
		folderPath, _ := filepath.Abs(filepath.Join(root, folderName))
		jobFolder[i].uri = folderPath
		os.MkdirAll(folderPath, os.ModePerm)
		splatJobStatus(folderPath, f.status)
		if f.tpl != nil {
			tplPath, err := filepath.Abs(filepath.Join("template", *f.tpl))
			if err != nil {
				slog.Error(fmt.Sprintf("cannot ABS template: %s", err))
				os.Exit(1)
			}
			dst := folderPath
			err = dCopy.Copy(tplPath, dst)
			if err != nil {
				slog.Error(fmt.Sprintf("cannot copy template: %s", err))
			}
		}
	}

	//display status of lock files
	jobs.ParseJobs()
	slog.Info("--- resetting status lock files ---------------------------------------------")

	//display all the status.lock files (and any extras)
	for _, j := range jobs.Known {
		extra := true
		for _, f := range jobFolder {
			if f.uri == string(j.Id) {
				slog.Debug(fmt.Sprintf("job %s stats=%s", j.Id, f.status))
				extra = false
			}
		}
		if extra {
			slog.Debug(fmt.Sprintf("job %s -----ADDITIONAL FOLDER-----", j.Id))
		}
	}
}
