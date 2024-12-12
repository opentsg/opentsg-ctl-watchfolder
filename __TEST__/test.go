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

	"github.com/phsym/console-slog"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/job"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/log"
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
}

var folders = []initState{
	initState{id: 123, status: "NEW"},
	initState{id: 124, status: "QUEUED"},
	initState{id: 125, status: "NEW test"},
	initState{id: 126, status: "FAILED"},
	initState{id: 127, status: "RUNNING"},
	initState{id: 155, status: "QUEUED"},
	initState{id: 157, status: "CANCELLED"},
	initState{id: 1538, status: "COMPLETED"},
}

var root = "./jobs"

func splatJobStatus(j job.JobInfo, status string) {
	path := filepath.Join(string(j.Id), "_status.lock")
	s := []byte(status)

	if err := os.WriteFile(path, s, 0644); err != nil {
		slog.Error(fmt.Sprintf("cannot write to %s", j.Id))
		return
	}
}

func main() {
	log.UsePrettyDebugLogger()
	slog.Info("__TEST__ reset job stats")

	//make all the folders
	for i, j := range folders {
		folderName := fmt.Sprintf("job%04d", j.id)
		folderPath, _ := filepath.Abs(filepath.Join(root, folderName))
		folders[i].uri = folderPath
		os.MkdirAll(folderPath, os.ModePerm)
	}

	//display status of lock files
	jobs := &job.JobManagement{
		Folder:       root,
		LockFileName: "_status.lock",
	}
	jobs.ParseJobs()
	slog.Info("--- resetting status lock files ---------------------------------------------")

	//write all the status.lock files
	for _, j := range jobs.Known {
		extra := true
		for _, f := range folders {
			//update the status if the files have the same abs path
			if f.uri == string(j.Id) {
				splatJobStatus(j, f.status)
				slog.Debug(fmt.Sprintf("job %s stats=%s", j.Id, f.status))
				extra = false
			}
		}
		if extra {
			slog.Debug(fmt.Sprintf("job %s -----ADDITIONAL FOLDER-----", j.Id))
		}
	}
}
