// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var tsgApp = "opentsg-node"
var optVersion = "-version"

func (j *JobInfo) getVersion() (version string, errMsg string) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	//run the app & capture stdout
	cmd := exec.Command(tsgApp, optVersion)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	cmd.Run()

	c := " \n\r\t"
	return strings.Trim(outBuf.String(), c), strings.Trim(errBuf.String(), c)
}

func (j *JobInfo) jobNodeLogger(buf io.Writer, rc io.ReadCloser,
	jobs *JobManagement, start int64) {
	// jLog, jobFile := log.JobLogger(string(j.jobLogPath))
	// defer jobFile.Close()

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Info(line)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("cannot redirect output to file", "err", err)
	}
	//clear the job running interlock
	jobs.JobRunning = nil

	// log the end times
	end := time.Now().UnixMilli()
	j.ActualEndDate = j.TimeStamp()
	j.ActualDuration = int(end - start)
	slog.Info(fmt.Sprintf("ending job at %s", j.ActualEndDate))
	slog.Info(fmt.Sprintf("duration of %d ms", j.ActualDuration))

	j.SetJobStatus(COMPLETED, "")

}

func (j *JobInfo) runJob(jobs *JobManagement) error {
	_dbg := fmt.Sprintf("         |%04d end-check", j.XjobId)
	var outBuf bytes.Buffer
	// var errBuf bytes.Buffer

	// make a logger for the user's job and close the file handle when done
	// jLog, jobFile := log.JobLogger(string(j.jobLogPath))
	// defer jobFile.Close()

	//setup the command to run
	mainJson := filepath.Join(string(j.XfolderPath), "main.json")
	// optRun := fmt.Sprintf("-c %s -output %s -log stdout -debug", mainJson, string(j.folderPath))
	argRun := []string{
		"-c", mainJson,
		"-jobid", j.IdString(),
		"-output", string(j.XfolderPath),
		"-log", "stdout",
		"-debug",
	}
	cmd := exec.Command(tsgApp, argRun...)

	j.Xcli = tsgApp + " " + strings.Join(argRun, " ")
	jobs.Xcli = j.Xcli

	// pipe stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	start := time.Now().UnixMilli()
	j.ActualStartDate = j.TimeStamp()
	slog.Info(fmt.Sprintf("% start @ %s", _dbg, j.ActualStartDate))

	// start the logging goroutine - tee to stdout and the job's log file
	go j.jobNodeLogger(&outBuf, stdout, jobs, start)

	//--------------------------------
	// run the node while logging to the log file and stdout
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	//--------------------------------

	return nil
}
