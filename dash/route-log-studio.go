// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

// package dash provides a simple dashboard for the job controller
func RouteStudioLogs(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	path, err := findStudioLogFilePath(jobId)
	if err != nil {
		// return an error view
		tpl.err.ExecuteTemplate(w, "page", TDErr{Title: jobId, Error: "Studio Logs not found"})
		return
	}
	slog.Debug("showing studio logs", "job", jobId, "studio-log", path)

	var j *job.JobInfo
	for i, jj := range jobs.Known {
		if jobId == jj.IdString() {
			j = &jobs.Known[i]
			break
		}
	}
	if j == nil {
		// return an error view
		tpl.err.ExecuteTemplate(w, "page", TDErr{Title: jobId, Error: "Job Logs not found"})
		return
	}
	logs := j.GetStudioLogs()
	data := TDStudioLogs{
		Title:     jobId + " logs (opentsg-studio)",
		Ptr:        "../", // relative path to the root folder ..
		L:         logs,
		LogSource: path,
	}

	err = tpl.studioLogs.ExecuteTemplate(w, "page", data)
	if err != nil {
		slog.Error("logs template render error", "err", err)
	}
}
