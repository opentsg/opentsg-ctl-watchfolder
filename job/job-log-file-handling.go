// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
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
	RunId       string    `json:"RunID"`
	WidgetId    string    `json:"WidgetID"`
	FrameNumber int       `json:"FrameNumber"`
}
type NodeLogLines struct {
	ErrorCount   int
	FrameCount   int
	FrameTotal   int
	RunCount     int
	LastModified time.Time
	LastError    string
	Lines        []NodeLogLine
}

// return a summary of the logs
func (j *JobInfo) GetNodeLogs() *NodeLogLines {
	path := string(j.XjobLogPath)
	logfileMeta, err := os.Stat(path)
	if err != nil {
		// no log file - just return
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		// maybe the file is locked - either way, return
		return nil
	}
	defer file.Close()

	ref := NodeLogLine{}
	logs := NodeLogLines{}
	logs.LastModified = logfileMeta.ModTime()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()
		logLine := NodeLogLine{}
		json.Unmarshal(line, &logLine)
		logs.Lines = append(logs.Lines, logLine)
		if ref.RunId != logLine.RunId {
			//if this is a new run then increment the run count and zero errors
			ref.RunId = logLine.RunId
			logs.RunCount += 1
			logs.ErrorCount = 0
		}
		if logLine.Level == "ERROR" {
			logs.ErrorCount += 1
			logs.LastError = logLine.Msg
		}
		logs.FrameCount = logLine.FrameNumber + 1
	}
	return &logs
}

// return lines from the studio logs
func (j *JobInfo) GetStudioLogs() *[]string {
	path := string(j.XstudioLogPath)
	_, err := os.Stat(path)
	if err != nil {
		// no log file - just return
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		// maybe the file is locked - either way, return
		return nil
	}
	defer file.Close()

	logs := []string{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := string(scanner.Bytes())
		logs = append(logs, line)
	}
	return &logs
}
