// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Check the job logs when the run stops
func (j *JobInfo) GetMostRecentLogs() map[string]interface{} {
	absFolder, _ := filepath.Abs(filepath.Join(string(j.XjobLogPath), fmt.Sprintf("job%04d", j.XjobId)))

	contents, err := os.ReadFile(absFolder)
	if err != nil {
		return nil
	}

	var logJson map[string]interface{}
	err = json.Unmarshal(contents, &logJson)
	if err != nil {
		return nil
	}
	return logJson
}

// Check the job logs when the run stops
func (j *JobInfo) JobEndCheck(jobs *JobManagement) {
	// _dbg := fmt.Sprintf("         |%04d end-check", j.XjobId)

	//we might already have checked this job - fast return
	if j.Xstate != StateRunning {
		return
	}
	//clear the checking flag
	j.Xstate = StateDone

	logs := j.GetNodeLogs()

	// if there is no log file, then something went wrong
	if logs == nil {
		j.SetJobStatus(FAILED, "no logs produced by render")
		return
	}

	if logs.errorCount > 0 {
		j.SetJobStatus(FAILED, logs.lastError)
	}
}
