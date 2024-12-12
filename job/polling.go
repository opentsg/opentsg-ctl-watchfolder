// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package job

import (
	"fmt"
	"log/slog"
	"time"
)

const interval = 5

type poller struct {
	ticker *time.Ticker // periodic ticker
	add    chan string  // new URL channel
}

func (jobs *JobManagement) StartPolling() *poller {
	rv := &poller{
		ticker: time.NewTicker(time.Second * interval),
		add:    make(chan string),
	}
	go rv.run(jobs)
	return rv
}

// run initiates an infinite poller to check all the folders
func (p *poller) run(jobs *JobManagement) {
	defer jobs.Wg.Done()
	for {
		select {
		case <-p.ticker.C:
			// When the ticker fires parse and handle all the jobs
			slog.Debug(fmt.Sprintf("=======  ===== %s @ %s", time.Now().Format(time.DateTime), jobs.Folder))

			jobs.ParseJobs()
			jobs.HandleJobs()
		case u := <-p.add:
			// At any time (other than when we're harvesting),
			// we can process a request to add a new URL
			slog.Info("Add another value ...", "value", u)
		}
	}
}
