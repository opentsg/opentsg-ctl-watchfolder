// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

// package dash provides a simple dashboard for the job controller

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

var jobs *job.JobManagement

func ShowDashboard(port int, eFs embed.FS, jobsToView *job.JobManagement, isProductionLogging bool) {
	jobs = jobsToView
	initTemplates(eFs)
	r := chi.NewRouter()

	if !isProductionLogging {
		// use the default logger when not in production mode
		r.Use(middleware.Logger)
	}
	// recover from panics and set return status
	r.Use(middleware.Recoverer)

	//set up routes
	r.Get("/", RouteJobs)
	r.Get("/{jobId}", RouteNodeLogs)
	r.Get("/{jobId}/studio-logs", RouteStudioLogs)

	// simple embedded file server for logs & static images, pages etc.
	embedFileServer(r, eFs, "/r/", "www")
	listenAddr := fmt.Sprintf("%s:%d", "", port)
	slog.Info(fmt.Sprintf("Listening on port %d", port))

	// run the server in a thread
	go http.ListenAndServe(listenAddr, r)
}
