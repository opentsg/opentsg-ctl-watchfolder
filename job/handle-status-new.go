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
)

// handle responses for a new job
func (j *JobInfo) HandleNewJob() {
	_dbg := fmt.Sprintf("         |%04d new", j.jobId)
	slog.Debug(fmt.Sprintf("%s (%s)", _dbg, j.meta))
	switch j.meta {
	case "test", "alive_response":
		msg := fmt.Sprintf("job%04d NEW active controller check", j.jobId)
		// check that the executable runs - return the version
		version, errMsg := j.getVersion()
		if len(errMsg) == 0 {
			j.SetJobStatus(NEW, version)
			slog.Info(msg, "job", j.Id, "res", j.meta)
		} else {
			slog.Error(msg, "job", j.Id, "res", j.meta)
		}
	}
}
