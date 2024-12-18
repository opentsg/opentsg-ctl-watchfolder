// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

// package dash provides a simple dashboard for the job controller

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/mrmxf/opentsg-ctl-watchfolder/job"
)

func ShowDashboard(port int, eFs embed.FS, jobs *job.JobManagement) {
	initTemplates(eFs)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Go to <a href=\"dash/\">dash/</a>!"))
	// })
	r.Get("/dash", JobsPageHandler)
	embedFileServer(r, eFs, "/", "www")
	listenAddr := fmt.Sprintf("%s:%d", "", port)
	http.ListenAndServe(listenAddr, r)
}
