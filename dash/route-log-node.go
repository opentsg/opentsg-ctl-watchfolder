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

type TplLogs struct {
	L         *job.NodeLogLines
	Title     string
	LogSource string
}

// locate the Node logfile
func findNodeLogFilePath(id string) (logFilePath string, err error) {
	logFilePath = filepath.Join(jobs.Folder, id, "_logs", id+".log")
	_, err = os.Stat(logFilePath)
	return
}

// locate the Studio logfile
func findStudioLogFilePath(id string) (logFilePath string, err error) {
	logFilePath = filepath.Join(jobs.Folder, id, "opentsg-studio")
	_, err = os.Stat(logFilePath)
	return
}

// package dash provides a simple dashboard for the job controller
func RouteShowLogs(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	nodeLogPath, err := findNodeLogFilePath(jobId)
	if err != nil {
		// assemble main
		err = tpl["main"].Execute(w, TplMain{
			Title: "job not found",
			Main:  "Job Logs not found",
		})
		if err != nil {
			slog.Error("main template render error", "err", err)
		}
		return
	}
	studioLogPath, _ := findStudioLogFilePath(jobId)
	slog.Info("showing logs", "job", jobId, "log", nodeLogPath, "studio", studioLogPath)

	var j *job.JobInfo
	for i, jj := range jobs.Known {
		if jobId == jj.IdString() {
			j = &jobs.Known[i]
			break
		}
	}
	if j == nil {
		// assemble main
		err = tpl["main"].Execute(w, TplMain{
			Title: "job not found",
			Main:  "Job Logs not found",
		})
	}
	data := TplLogs{
		Title:     jobId + " logs (opentsg-node)",
		L:         j.GetNodeLogs(),
		LogSource: nodeLogPath,
	}
	err = dashTpl.logs.Execute(w, data)
	if err != nil {
		slog.Error("logs template render error", "err", err)
	}
}
