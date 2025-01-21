// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

type NodeLogLine struct {
	Time        time.Time `json:"time"`
	Level       string    `json:"level"`
	Msg         string    `json:"msg"`
	StatusCode  string    `json:"StatusCode"`
	RunID       string    `json:"RunID"`
	WidgetID    string    `json:"WidgetID"`
	FrameNumber int       `json:"FrameNumber"`
}
type NodeLogLines struct {
	errorCount   int
	frameCount   int
	frameTotal   int
	runCount     int
	lastModified time.Time
	lastError    string
	lines        []NodeLogLine
}

// return a summary of the logs
func (j *JobInfo) GetNodeLogs() *NodeLogLines {
	logfileMeta, err := os.Stat(string(j.XjobLogPath))
	if err != nil {
		// no log file - just return
		return nil
	}

	file, err := os.Open(string(j.XjobLogPath))
	if err != nil {
		// maybe the file is locked - either way, return
		return nil
	}
	defer file.Close()

	ref := NodeLogLine{}
	logs := NodeLogLines{}
	logs.lastModified = logfileMeta.ModTime()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()
		logLine := NodeLogLine{}
		json.Unmarshal(line, &logLine)
		logs.lines = append(logs.lines, logLine)
		if ref.RunID != logLine.RunID {
			//if this is a new run then increment the run count and zero errors
			ref.RunID = logLine.RunID
			logs.runCount += 1
			logs.errorCount = 0
		}
		if logLine.Level == "ERROR" {
			logs.errorCount += 1
			logs.lastError = logLine.Msg
		}
		logs.frameCount = logLine.FrameNumber + 1
	}
	return &logs
}
