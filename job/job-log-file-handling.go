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
	FrameNumber string    `json:"FrameNumber"`
}
type NodeLogLines struct {
	errorCount   int
	frameCount   int
	runCount     int
	lastModified time.Time
	lines        []NodeLogLine
}

func (j *JobInfo) GetNodeLogs(status JobStatusEnum, meta string) *NodeLogLines {
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
	logLine.lastModified = logfileMeta.ModTime()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()
		logLine := NodeLogLine{}
		json.Unmarshal(line, &logLine)
	}
	return &logs
}
