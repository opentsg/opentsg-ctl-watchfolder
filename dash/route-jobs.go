// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sort"
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
	// need to reverse sort the jobs array:
	known := jobs.Known
	sort.Slice(known, func(i, j int) bool {
		return known[i].XjobId > known[j].XjobId
	})
	for i, j := range known {
		//format the duration into the Xage field
		known[i].Xage = "-"
		if j.ActualDuration > 1 {
			known[i].Xage = fmt.Sprintf("%d secs", j.ActualDuration/1000)
		}
		if j.ActualDuration > 120000 {
			known[i].Xage = fmt.Sprintf("%d mins", j.ActualDuration/60000)
		}
		if j.ActualDuration > (1000 * 60 * 60 * 3) {
			known[i].Xage = fmt.Sprintf("%d hrs", j.ActualDuration/(1000*60*60))
		}
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
