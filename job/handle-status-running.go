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

// QueuePosition returns the integer position of the job in the queue or error
func (j *JobInfo) JobCompletion(jobs *JobManagement) (int, int, int, error) {
	return 50, 11, 22, nil
}

// feed back progress to the UX
// Check that this is the actual running job. If not, then it failed somehow
func (j *JobInfo) RunningJob(jobs *JobManagement) {
	_dbg := fmt.Sprintf("         |%04d running", j.XjobId)

	if jobs.JobRunning == nil || jobs.JobRunning.Id != j.Id {
		// this should never happen! something got stuck and maybe the controller
		// or the runner crashed - either way mark it as failed and the user will
		// reset the job (or raise a support ticket!)
		meta := "unknown error - probably opentsg-node crash?"
		j.SetJobStatus(FAILED, meta)
		j.Xstate = StateDone
		slog.Debug(fmt.Sprintf("%s %s", _dbg, meta))
		return
	}

	progress, frame, total, err := j.JobCompletion(jobs)

	// update the status of the running job
	if err == nil {
		meta := fmt.Sprintf("%d %d %d", progress, frame, total)
		j.SetJobStatus(RUNNING, meta)
		slog.Debug(fmt.Sprintf("%s %s", _dbg, meta))
	}
}
