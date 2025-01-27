package dash

import (
	"embed"
	"html/template"

	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

// Define some data template data structs
type TDErr struct {
	Title string
	Error template.HTML
}

type TDJobsMain struct {
	Title	          string
	Folder          string
	JobCount        int
	QueueDepth      int
	JobRunningIdent string
	JobCli          string
	List            *[]TDJob
}

type TDJob struct {
	J         job.JobInfo
	NodeLog   string
	StudioLog string
}

type TDNodeLogs struct {
	Title     string
	LogSource string
	L         *job.NodeLogLines
}

type TDStudioLogs struct {
	Title     string
	LogSource string
	L         *[]string
}

// This struct is for managing all the templates used for rendering. A template
// is a set of named templates loaded as a group. The xxxSrc array holds the
// source filename in the embedded file system. the xxx template pointer holds
// the parsed group of templates. When executing, use the name "page" to ref
// the generic outer template
type DashTemplates struct {
	err           *template.Template
	errSrc        []string
	jobsMain      *template.Template
	jobsMainSrc   []string
	nodeLogs      *template.Template
	nodeLogsSrc   []string
	studioLogs    *template.Template
	studioLogsSrc []string
}

var tpl = DashTemplates{
	errSrc: []string{
		"dash/templates/page.html",
		"dash/templates/main-err.html",
	},
	jobsMainSrc: []string{
		"dash/templates/page.html",
		"dash/templates/header-refresh.html",
		"dash/templates/main-jobs-table.html",
		"dash/templates/job-list.html",
	},
	nodeLogsSrc: []string{
		"dash/templates/page.html",
		"dash/templates/head.html",
		"dash/templates/main-logs-node.html",
	},
	studioLogsSrc: []string{
		"dash/templates/page.html",
		"dash/templates/head.html",
		"dash/templates/main-logs-studio.html",
	},
}

// iterate over the embedded source files and init the templates
func initTemplates(fs embed.FS) {
	// make a new logs template by parsing the collection of named templates
	tpl.err = template.Must(template.ParseFS(fs, tpl.errSrc...))
	tpl.jobsMain = template.Must(template.ParseFS(fs, tpl.jobsMainSrc...))
	tpl.nodeLogs = template.Must(template.ParseFS(fs, tpl.nodeLogsSrc...))
	tpl.studioLogs = template.Must(template.ParseFS(fs, tpl.studioLogsSrc...))
}
