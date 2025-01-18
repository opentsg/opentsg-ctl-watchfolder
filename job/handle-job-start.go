// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"fmt"
	"log/slog"

	"github.com/opentsg/opentsg-ctl-watchfolder/shell"
)

var dummyJob = "sleep 2 && echo 'dummy job done'"

func (j *JobInfo) StartJob() {
	j.SetJobStatus(RUNNING, "")
	slog.Debug(fmt.Sprintf("job%04d  start", j.XjobId))
	res := shell.StreamShellSnippet(dummyJob, nil)
	if res == nil {
		j.SetJobStatus(FAILED, "")
		slog.Debug(fmt.Sprintf("job%04d  failed", j.XjobId))
	} else {

		j.SetJobStatus(RUNNING, "")
		slog.Debug(fmt.Sprintf("job%04d  running", j.XjobId))
	}
}
