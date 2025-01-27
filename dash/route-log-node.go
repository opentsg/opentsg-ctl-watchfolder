// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

// locate the Node log file
func findNodeLogFilePath(id string) (logFilePath string, err error) {
	logFilePath = filepath.Join(jobs.Folder, id, jobs.LogsFolder, id+".log")
	_, err = os.Stat(logFilePath)
	if err != nil {
		return "", err
	}
	return
}

// locate the Studio log file
func findStudioLogFilePath(id string) (logFilePath string, err error) {
	logFilePath = filepath.Join(jobs.Folder, id, jobs.LogsFolder, jobs.LogStudioName)
	_, err = os.Stat(logFilePath)
	if err != nil {
		return "", err
	}
	return
}

// package dash provides a simple dashboard for the job controller
func RouteNodeLogs(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	path, err := findNodeLogFilePath(jobId)
	if err != nil {
		// return an error view
		tpl.err.ExecuteTemplate(w, "page", TDErr{Title: jobId, Error: "Job Logs not found"})
		return
	}
	slog.Debug("showing node logs", "job", jobId, "log", path)

	var j *job.JobInfo
	for i, jj := range jobs.Known {
		if jobId == jj.IdString() {
			j = &jobs.Known[i]
			break
		}
	}
	if j == nil {
		// return an error view
		tpl.err.ExecuteTemplate(w, "page", TDErr{Title: jobId, Error: "Job not found"})
		return
	}
	logs := j.GetNodeLogs()
	data := TDNodeLogs{
		Title:     jobId + " logs (opentsg-node)",
		Ptr:        "../", // relative path to the root folder ..
		L:         logs,
		LogSource: path,
	}
	err = tpl.nodeLogs.ExecuteTemplate(w, "page", data)
	if err != nil {
		slog.Error("logs template render error", "err", err)
	}
}
