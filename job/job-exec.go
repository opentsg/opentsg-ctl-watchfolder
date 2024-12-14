// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/log"
)

var tsgApp = "msgtsg-node"
var optVersion = "--version"

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

func jobNodeLogger(wg *sync.WaitGroup, buf io.Writer, rc io.ReadCloser, name string) {
	jLog, jobFile := log.JobLogger(string(j.jobLogPath))
	defer jobFile.Close()
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		jLog.Info(name, line)
		fmt.Fprintln(buf, line)
	}
	if err := scanner.Err(); err != nil {
		jLog.Error("cannot redirect output to file", "err", err)
	}
	wg.Done()
}

func (j *JobInfo) runJob(jobs *JobManagement) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	// make a logger for the user's job and close the file handle when done
	jLog, jobFile := log.JobLogger(string(j.jobLogPath))
	defer jobFile.Close()
	//run the app & capture stdout
	cmd := exec.Command(tsgApp, optVersion)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	start := time.Now().UnixMilli()
	j.ActualStartDate = j.TimeStamp()
	jLog.Info(fmt.Sprintf("starting job at %s", j.ActualStartDate))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}

	j.Wg.Add(1)
	// start the logging goroutine - tee to stdout and the job's log file
	go jobNodeLogger(&wg, &outBuf, stdout, "stdout")
	// run the node while logging to the log file and stdout
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "command failed")
	}
	// await  the node's stdout closing
	wg.Wait()

	// log the end times
	end := time.Now().UnixMilli()
	j.ActualEndDate = j.TimeStamp()
	j.ActualDuration = int(end - start)
	jLog.Info(fmt.Sprintf("ending job at %s", j.ActualEndDate))
	jLog.Info(fmt.Sprintf("duration of %d ms", j.ActualDuration))

	//clear the running job lock
	jobs.JobRunning = nil
	j.SetJobStatus(COMPLETED, "")
}
