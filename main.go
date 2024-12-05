// mk8 simple job scheduler
package main

import (
	"bufio"
	"fmt"
	"github.com/phsym/console-slog"
	"gitlab.com/mrmxf/opentsg-ctl-mk8/shell"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const jobsFolder = "jobs/"
const jobGlob = "job*"
const jobLockfile = "_status.lock"

var jobInProgress = 0

type JobStatusType int

const (
	JobUNKNOWN   JobStatusType = iota
	JobCREATED   JobStatusType = iota
	JobINVALID   JobStatusType = iota
	JobSUBMITTED JobStatusType = iota
	JobQUEUED    JobStatusType = iota
	JobSTART     JobStatusType = iota
	JobRUNNING   JobStatusType = iota
	JobFAILED    JobStatusType = iota
)

var jobStatusString = []string{
	"UNKNOWN",
	"CREATED",
	"INVALID",
	"SUBMITTED",
	"QUEUED",
	"START",
	"RUNNING",
	"FAILED",
}

type JobInfo struct {
	folderPath string
	id         int
	status     JobStatusType
	statusStr  string
}

func setupLogger() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)
}

func getJobStatus(folderPath string) (int, string, error) {
	path := filepath.Join(folderPath, jobLockfile)
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		slog.Debug(fmt.Sprintf("   err %s cannot be opened", path))
		return 0, "", err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		r, _ := regexp.Compile(`([0-9]{1,2})\s([^\s]+).*`)
		textLine := scanner.Text()
		slog.Debug(fmt.Sprintf("status %s = \"%s\"", path, textLine))
		match := r.FindStringSubmatch(textLine)
		if len(match) == 3 {
			// e.g. meta[short] = "This is short help"
			status, _ := strconv.Atoi(match[1])
			slog.Debug(fmt.Sprintf(" parse %s = (%d,%s)", path, status, match[2]))
			return status, match[2], nil
		}
	}
	return 0, "0000", nil
}

func getJobStats() []JobInfo {
	//look for all shell scripts in the clogrc folder
	searchFolder := jobsFolder + jobGlob
	slog.Debug("  glob " + searchFolder)

	jobsRaw, _ := filepath.Glob(searchFolder)
	slog.Debug(fmt.Sprintf(" found %d folders", len(jobsRaw)))
	jobs := []JobInfo{}
	//add each script found
	for _, j := range jobsRaw {
		slog.Debug(fmt.Sprintf(" check %s ", j))
		status, statusStr, _ := getJobStatus(j)
		id, _ := strconv.Atoi(j[len(j)-4:])
		job := JobInfo{
			folderPath: j,
			id:         id,
			statusStr:  statusStr,
			status:     JobStatusType(status),
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func setJobStatus(j JobInfo, status JobStatusType, meta string) {
	path := filepath.Join(j.folderPath, jobLockfile)
	i := int(status)
	s := []byte(fmt.Sprintf("%d %s %s", i, jobStatusString[i], meta))
	if err := os.WriteFile(path, s, 0644); err != nil {
		slog.Error(fmt.Sprintf("job%04d  status update failed - cannot write to %s", j.folderPath))
		return
	}
}
func queueJob(j JobInfo) {
	setJobStatus(j, JobQUEUED, "")
	slog.Debug(fmt.Sprintf("job%04d  queue", j.id))
}

func startJob(j JobInfo) {
	setJobStatus(j, JobSTART, "")
	slog.Debug(fmt.Sprintf("job%04d  start", j.id))
	res := shell.StreamShellSnippet("sleep 0.5", nil)
	if res == nil {
		setJobStatus(j, JobFAILED, "Just testing")
		slog.Debug(fmt.Sprintf("job%04d  failed", j.id))
	} else {

		setJobStatus(j, JobRUNNING, "")
		slog.Debug(fmt.Sprintf("job%04d  running", j.id))
	}
}

func handleJob(j JobInfo) {
	slog.Debug(fmt.Sprintf("job%04d  -------", j.id))
	switch j.status {
	case JobSUBMITTED:
		queueJob(j)
	case JobQUEUED:
		if jobInProgress == 0 {
			startJob(j)
		}
	case JobSTART:

	default:
		return
	}
}

func main() {
	setupLogger()
	slog.Info("Minikube Job Controller started")
	jobs := getJobStats()
	for _, j := range jobs {
		handleJob(j)
	}
}
