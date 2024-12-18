// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package dash

import "net/http"

// package dash provides a simple dashboard for the job controller
func JobsPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Jobs page"))
}
