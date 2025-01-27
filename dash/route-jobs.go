// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"
)

const (
	Duration_2min = (1000 * 60 * 2)
	Duration_3hr  = (1000 * 60 * 60 * 3)
	Duration_48hr = (1000 * 60 * 60 * 48)
)

func friendlyDuration(durMs int) string {
	if durMs == 0 {
		return "-"
	}
	if durMs <= Duration_2min {
		return fmt.Sprintf("%d secs", durMs/1000)
	}
	if durMs <= Duration_3hr {
		return fmt.Sprintf("%d mins", durMs/(1000*60))
	}
	if durMs <= Duration_48hr {
		return fmt.Sprintf("%d hrs", durMs/(1000*60*60))
	}
	return fmt.Sprintf("%d days", durMs/(1000*60*60*24))
}

// package dash provides a simple dashboard for the job controller
func RouteJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobList := []TDJob{}
	jobsData := TDJobsMain{
		Title:      "Jobs (opentsg-studio)",
		Ptr:        "", // relative path to the root folder nothing
		Folder:     jobs.Folder,
		JobCount:   len(jobs.Known),
		QueueDepth: len(jobs.Queue),
		List:       &jobList,
	}
	var err error

	// need to reverse sort the jobs array:
	known := jobs.Known
	sort.Slice(known, func(i, j int) bool {
		return known[i].XjobId > known[j].XjobId
	})
	//build the job data - newest first
	for i, j := range known {
		//format the duration & age
		known[i].XDurationStr = friendlyDuration(j.ActualDuration)
		start, err := time.Parse("2006-01-02 15:04:05", j.ActualStartDate)
		if err != nil {
			known[i].XAgeStr = ""
		} else {
			known[i].XAgeStr = friendlyDuration(int(time.Since(start) / 1000000))
		}

		nodeLogPath, _ := findNodeLogFilePath(j.IdString())
		studioLogPath, _ := findStudioLogFilePath(j.IdString())
		jobList = append(jobList, TDJob{
			J:         j,
			NodeLog:   nodeLogPath,
			StudioLog: studioLogPath,
		})
	}

	if jobs.JobRunning != nil {
		jobsData.JobRunningIdent = fmt.Sprintf("job%04d", jobs.JobRunning.XjobId)
		jobsData.JobCli = jobs.Xcli
	}
	//render the Jobs
	err = tpl.jobsMain.ExecuteTemplate(w, "page", jobsData)
	if err != nil {
		slog.Error("jobsMain template render error", "err", err)
	}

}
