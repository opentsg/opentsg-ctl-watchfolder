// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

// package job contains the basic processing for jobs.
//
// There is an interface file for simple minikube io, a definitions file and a
// logic file that defines rules

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

// ParseLockFileString parses the plain text lockfile and returns values.
func (j *JobInfo) ParseLockFileString(s string) (status string, meta string, err error) {

	// this regexp matches "NEW test", "RUNNING 70%" etc
	r, _ := regexp.Compile(`(\w+)\s*(.*)`)

	match := r.FindStringSubmatch(s)
	if len(match) != 3 {
		return "", "", errors.New("invalid lockfile string")
	}
	c := " \n\r\t"
	return strings.Trim(match[1], c), strings.Trim(match[2], c), nil
}

// GetJobMetadata retrieves the metadata from job folders
func (j *JobInfo) ReadLockFileMetadata() (status string, meta string, err error) {

	file, err := os.Open(string(j.lockfilePath))
	if err != nil {
		slog.Debug(fmt.Sprintf("   err %s cannot be opened", j.lockfilePath))
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	// just use the first line (for now)
	textLine := scanner.Text()
	status, meta, err = j.ParseLockFileString(textLine)
	return
}

func (j *JobInfo) SetJobStatus(status JobStatusEnum, meta string) {
	//set and write to lock file
	s := []byte(fmt.Sprintf("%s %s", status, meta))
	err := os.WriteFile(string(j.lockfilePath), s, 0644)

	if err != nil {
		slog.Error(fmt.Sprintf("job%04d  status update failed - cannot write to %s", j.jobId, string(j.folderPath)))
		return
	}

	// check by reading back
	st, me, _ := j.ReadLockFileMetadata()
	j.Status = JobStatusEnum(st)
	j.meta = me
}
