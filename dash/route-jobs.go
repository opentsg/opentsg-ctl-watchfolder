// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

// package dash provides a simple dashboard for the job controller
func RouteJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobsHTML := bytes.Buffer{}
	jobsData := TplJobs{
		Folder:     jobs.Folder,
		JobCount:   len(jobs.Known),
		QueueDepth: len(jobs.Queue),
	}
	var err error

	//render each job outside the template to allow sorting etc.
	jobsData.JobTableHTML = ""
	for _, j := range jobs.Known {
		tmp := bytes.Buffer{}
		err = tpl["job"].Execute(&tmp, j)
		if err != nil {
			slog.Error("job template render error", "job", j.XjobId, "err", err)
		}
		jobsData.JobTableHTML += template.HTML(tmp.Bytes())
	}
	if jobs.JobRunning != nil {
		jobsData.JobRunningIdent = fmt.Sprintf("job%04d", jobs.JobRunning.XjobId)
		jobsData.JobCli = jobs.Xcli
	}
	//render the Jobs
	err = tpl["jobs"].Execute(&jobsHTML, jobsData)
	if err != nil {
		slog.Error("jobs template render error", "err", err)
	}

	// assemble main
	err = tpl["main"].Execute(w, TplMain{
		Title: "Test render",
		Main:  template.HTML(jobsHTML.Bytes()),
	})
	if err != nil {
		slog.Error("main template render error", "err", err)
	}

}
