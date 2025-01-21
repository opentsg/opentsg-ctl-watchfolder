package main

import (
	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

var j = job.JobInfo{
	//internal properties
	XjobId:        100,
	XfolderPath:   ".",
	XlockFilePath: "_status.lock",
	XjobLogPath:   "job0003.log",
	Xmeta:         "",
	Xstate:        job.StateUnknown,
	XfirstSeenAt:  "2025-01-21 23:45:01",
	XqueuedAt:     "2025-01-21 23:45:02",
	Xcli:          "",
	XDurationStr:  "",
	XAgeStr:       "",

	//the following parameters are external and follow SMPTE ST2126:2020
	Id:              job.URL(""),           // URL pointing to the job instance in the job processor
	Type:            "render",              // Indicates the job type
	Profile:         job.URL(""),           //URL pointing to the job profile used by the job
	ProfileName:     "",                    //Name of the job profile used by the job.
	Execution:       job.URL(""),           //URL pointing to the jobExecution instance in the job processor
	Assignment:      job.URL(""),           // URL pointing to the jobAssignment instance in the executing service
	Input:           "",                    // Collection of input parameters that were provided in the job when it was created
	Status:          job.RUNNING,           //Status of the job
	Error:           nil,                   // Detailed info about the problem which caused the job. nil when not failed
	ActualStartDate: "2025-01-21 23:45:10", // Date in ISO 8601 format when job was queued for processing
	ActualEndDate:   "2025-01-21 23:45:15", //Date in ISO 8601 format when job completed, failed or canceled
	ActualDuration:  5123,                  //Job duration in milliseconds
	Output:          job.OutputInfo{},      // Collection of output results of the job that was executed

}

func main() {
	j.GetNodeLogs()
}
