// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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

func (j *JobInfo) runJob(jobs *JobManagement) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	jlog := log.JobLogger(string(j.joblogPath))
	//run the app & capture stdout
	cmd := exec.Command(tsgApp, optVersion)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	start := time.Now().UnixMilli()
	j.ActualStartDate = j.TimeStamp()
	jlog.info(fmt.Sprintf("starting job at %s", j.ActualStartDate))
	cmd.Run()
	end := time.Now().UnixMilli()
	j.ActualEndDate = j.TimeStamp()
	j.ActualDuration = int(end - start)
	jlog.info(fmt.Sprintf("ending job at %s", j.ActualEndDate))
	jlog.info(fmt.Sprintf("duration of %d ms", j.ActualDuration))

	//clear the running job lock
	jobs.JobRunning = nil
	j.SetJobStatus(COMPLETED, "")
}
