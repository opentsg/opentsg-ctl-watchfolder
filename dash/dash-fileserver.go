// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// embedFileServer conveniently sets up a http.embedFileServer handler to serve
// static files from a http.FileSystem.
func embedFileServer(r chi.Router, eFs embed.FS, route string, eFsRootPath string) {
	if strings.ContainsAny(route, "{}*") {
		slog.Error("embedFileServer mount route does not permit any URL parameters.")
		return
	}

	fSys, err := fs.Sub(eFs, eFsRootPath)
	if err != nil {
		slog.Error("embedFileServer cannot find embedded files")
		return
	}

	// check for trailing slash
	if route != "/" && route[len(route)-1] != '/' {
		r.Get(route, http.RedirectHandler(route+"/", http.StatusMovedPermanently).ServeHTTP)
		route += "/"
	}
	route += "*"

	r.Get(route,
		func(w http.ResponseWriter, r *http.Request) {
			rCtx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rCtx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(fSys)))
			fs.ServeHTTP(w, r)
		})
}
