// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"errors"
	"fmt"
	"log/slog"
)

// QueuePosition returns the integer position of the job in the queue or error
func (j *JobInfo) QueuePosition(jobs *JobManagement) (int, error) {
	for i, q := range jobs.Queue {
		if j.Id == q.Id {
			//position #1 in the queue means next to be scheduled
			return len(jobs.Queue) - i, nil
		}
	}
	return 0, errors.New("job not found in queue")
}

func (j *JobInfo) QueueJob(jobs *JobManagement) {
	_dbg := fmt.Sprintf("         |%04d queue", j.XjobId)
	pos, err := j.QueuePosition(jobs)

	// if job not in queue - add it
	if err != nil {
		slog.Debug(fmt.Sprintf("%s add", _dbg))
		jobs.Queue = append(jobs.Queue, j)
		pos, _ = j.QueuePosition(jobs)
	}

	meta := fmt.Sprintf("%d", pos)
	j.SetJobStatus(QUEUED, meta)
	j.XqueuedAt = j.TimeStamp()
	slog.Debug(fmt.Sprintf("%s pos %s", _dbg, meta))

	if pos == 1 && jobs.JobRunning == nil {
		// lock this job as the running job & remove from queue
		jobs.JobRunning = j
		jobs.Queue = jobs.Queue[1:]
		j.SetJobStatus(RUNNING, "")
		j.Xstate = StateRunning
		j.runJob(jobs)
	}
}
