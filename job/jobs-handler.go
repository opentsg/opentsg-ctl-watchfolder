// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
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
	"os"
	"path/filepath"
	"strconv"
)

const DEBUG_PARSER = false

// ParseJobs searches the main job folder for all matching job subfolders.
//
// Within the folder it looks for the lock file to communicate with the front
// end. It logs any errors it files to the logs folder.
func (jobs *JobManagement) ParseJobs() {
	globFolders := filepath.Join(jobs.Folder, globExpression)
	jobsRaw, _ := filepath.Glob(globFolders)
	if DEBUG_PARSER {
		slog.Debug(fmt.Sprintf(" found %d folders with glob(%s)", len(jobsRaw), globFolders))
	}

	// reset the seen status of each job
	for i, _ := range jobs.Known {
		jobs.Known[i].Xstate = StateUnknown
	}

	// look for lock file in each job folder found
	for _, jFolder := range jobsRaw {
		absFolder, _ := filepath.Abs(jFolder)
		tmp := JobInfo{
			XfolderPath:   URL(jFolder),
			XlockFilePath: URL(filepath.Join(jFolder, jobs.LockFileName)),
			XjobLogPath:   URL(filepath.Join(jFolder, jobs.JobLogName)),
			Id:            URL(absFolder),
		}
		status, meta, start, end, age, err := tmp.ReadLockFileMetadata()
		if err == nil {
			intId, _ := strconv.Atoi(jFolder[len(jFolder)-4:])
			tmp.XjobId = intId
			tmp.Status = JobStatusEnum(status)
			tmp.Xmeta = meta
			tmp.ActualStartDate = start
			tmp.ActualEndDate = end
			tmp.ActualDuration = age
			jobs.UpdateKnownJobs(&tmp)
			if DEBUG_PARSER {
				slog.Debug(fmt.Sprintf("status job%04d (%-9s) meta(%-12s)  << %s", tmp.XjobId, status, meta, tmp.XlockFilePath))
			}
		} else {
			// if the lockfile does not exist ignore error
			if !errors.Is(err, os.ErrNotExist) {
				slog.Debug(fmt.Sprintf("ERROR reading job%04d lockfile << %s", tmp.XjobId, tmp.XlockFilePath))
			}
		}
	}

	// set all unknown job states to deleted
	for i := range jobs.Known {
		if jobs.Known[i].Xstate == StateUnknown {
			jobs.Known[i].Xstate = StateDeleted
		}
	}

}

// iterate over the managed jobs and add new ones
func (jobs *JobManagement) UpdateKnownJobs(newJob *JobInfo) {
	for i, j := range jobs.Known {
		if j.Id == newJob.Id {
			//we've seen this job before - let's update from the lockfile
			jobs.Known[i].Status = newJob.Status
			jobs.Known[i].Xmeta = newJob.Xmeta
			//set the state to something benign if the job is unknown
			if jobs.Known[i].Xstate == StateUnknown {
				jobs.Known[i].Xstate = StateSeen
			}
			return
		}
	}
	//append the new job with the time we saw it
	newJob.XfirstSeenAt = newJob.TimeStamp()
	jobs.Known = append(jobs.Known, *newJob)
}

// iterate over the managed jobs and handle them. This happens once every
// polling cycle. The handlers update the jobs so the index of the array is
// needed when calling the handlers.
func (jobs *JobManagement) HandleJobs() {
	for i, j := range jobs.Known {
		_fmt := "job%04d  ----- %-12s -----------------------"
		_hdr := "============"
		switch j.Status {
		case NEW:
			slog.Debug(fmt.Sprintf(_fmt, j.XjobId, _hdr))
			jobs.Known[i].HandleNewJob()
		case QUEUED:
			slog.Debug(fmt.Sprintf(_fmt, j.XjobId, _hdr))
			jobs.Known[i].QueueJob(jobs)
		case RUNNING:
			slog.Debug(fmt.Sprintf(_fmt, j.XjobId, _hdr))
			jobs.Known[i].RunningJob(jobs)
		case COMPLETED:
			slog.Debug(fmt.Sprintf(_fmt, j.XjobId, _hdr))
			jobs.Known[i].JobEndCheck(jobs)
		case FAILED:
			slog.Debug(fmt.Sprintf(_fmt, j.XjobId, _hdr))
			jobs.Known[i].JobEndCheck(jobs)
		default:
			if j.Xstate == StateDeleted {
				slog.Debug(fmt.Sprintf(_fmt, j.XjobId, "Deleted"))
			} else {
				slog.Debug(fmt.Sprintf(_fmt, j.XjobId, j.Status))
			}
		}
	}
}
