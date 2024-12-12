// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"sync"
	"time"

	"github.com/longkai/rfc7807"
)

// a glob expression for matching job
const globExpression = "job????"
const jobLogFile = "_ctl-watchfolder.log"

// this is the status that is exported externally
type JobStatusEnum string

// slight superset of SMPTE ST 2126
const (
	NEW       JobStatusEnum = "NEW"
	QUEUED    JobStatusEnum = "QUEUED"
	RUNNING   JobStatusEnum = "RUNNING"
	COMPLETED JobStatusEnum = "COMPLETED"
	FAILED    JobStatusEnum = "FAILED"
	CANCELLED JobStatusEnum = "CANCELLED"
)

// this is the internal state tracking micro-tasks
type StateEnum int

const (
	StateNew               StateEnum = iota
	StateNewInitialised    StateEnum = iota
	StateQueuedSubmitted   StateEnum = iota
	StateQueuedWaiting     StateEnum = iota
	StateRunningInit       StateEnum = iota
	StateRunningOk         StateEnum = iota
	StateRunningWaiting    StateEnum = iota
	StateCompleted         StateEnum = iota
	StateFailed            StateEnum = iota
	StateFailedInvalid     StateEnum = iota
	StateCancelledInit     StateEnum = iota
	StateCancelledWaiting  StateEnum = iota
	StateCancelledComplete StateEnum = iota
	StateCancelledFailed   StateEnum = iota
)

var StatusMap = map[StateEnum]JobStatusEnum{
	StateNew:               NEW,
	StateNewInitialised:    NEW,
	StateQueuedSubmitted:   QUEUED,
	StateQueuedWaiting:     QUEUED,
	StateRunningInit:       RUNNING,
	StateRunningOk:         RUNNING,
	StateRunningWaiting:    RUNNING,
	StateCompleted:         COMPLETED,
	StateFailed:            FAILED,
	StateFailedInvalid:     FAILED,
	StateCancelledInit:     CANCELLED,
	StateCancelledWaiting:  CANCELLED,
	StateCancelledComplete: CANCELLED,
	StateCancelledFailed:   FAILED,
}

type LogLevelCode int
type LogLevelName string
type LogLevelMap map[LogLevelName]LogLevelCode

const (
	FATAL          LogLevelCode = 100
	ERROR          LogLevelCode = 200
	WARN           LogLevelCode = 300
	INFO           LogLevelCode = 400
	DEBUG          LogLevelCode = 500
	FUNCTION_START LogLevelCode = 450
	FUNCTION_END   LogLevelCode = 450
	JOB_START      LogLevelCode = 400
	JOB_UPDATE     LogLevelCode = 400
	JOB_END        LogLevelCode = 400
)

var LogLevel = LogLevelMap{
	"FATAL":          FATAL,
	"ERROR":          ERROR,
	"WARN":           WARN,
	"INFO":           INFO,
	"DEBUG":          DEBUG,
	"FUNCTION_START": FUNCTION_START,
	"FUNCTION_END":   FUNCTION_END,
	"JOB_START":      JOB_START,
	"JOB_UPDATE":     JOB_UPDATE,
	"JOB_END":        JOB_END,
}

type URL string

type OutputInfo struct {
	LogLocation URL
}

type ErrorInfo *rfc7807.ProblemDetail

type JobInfo struct {
	//internal properties
	jobId        int
	folderPath   URL
	lockfilePath URL
	joblogPath   URL
	meta         string
	state        StateEnum
	firstSeenAt  string
	queuedAt     string

	//the following parameters are external and follow SMPTE ST2126:2020
	// URL pointing to the job instance in the job processor
	Id URL
	// Indicates the job type
	Type string
	//URL pointing to the job profile used by the job
	Profile URL
	//Name of the job profile used by the job.
	ProfileName string
	//URL pointing to the jobExecution instance in the job processor
	Execution URL
	// URL pointing to the jobAssignment instance in the executing service
	Assignment URL
	// Collection of input parameters that were provided in the job when it was created
	Input string
	//Status of the job
	Status JobStatusEnum
	// Detailed info about the problem which caused the job. nil when not failed
	Error ErrorInfo
	// Date in ISO 8601 format when job was queued for processing
	ActualStartDate string
	//Date in ISO 8601 format when job completed, failed or canceled
	ActualEndDate string
	//Job duration in milliseconds
	ActualDuration int
	// Collection of output results of the job that was executed
	Output OutputInfo
}

type JobManagement struct {
	JobRunning   *JobInfo
	Known        []JobInfo
	Queue        []*JobInfo
	Folder       string
	LockFileName string
	JobLogName   string

	Wg sync.WaitGroup
}

// SMPTE ST2126 timestamp
func (j *JobInfo) TimeStamp() string {
	return time.Now().Format(time.DateTime)
}
