package dash

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/opentsg/opentsg-ctl-watchfolder/job"
)

type DashTemplateSrcFileMap map[string]string
type DashTemplateMap map[string]*template.Template

type TplMain struct {
	Title string
	Main  template.HTML
}

type TplJobs struct {
	Folder          string
	JobCount        int
	QueueDepth      int
	JobRunningIdent string
	JobCli          string

	JobTableHTML template.HTML
}

type TplJob struct {
	J         job.JobInfo
	NodeLog   string
	StudioLog string
}

var tplSrcFile = DashTemplateSrcFileMap{
	"main": "dash/templates/page-dash.html",
	"job":  "dash/templates/job-row.html",
	"jobs": "dash/templates/jobs-table.html",
}
var tpl = DashTemplateMap{}

type DashTemplates struct {
	jobs    *template.Template
	SrcJobs []string
	logs    *template.Template
	SrcLogs []string
}

var dashTpl = DashTemplates{
	SrcJobs: []string{
		"dash/templates/page-dash.html",
		"dash/templates/job-row.html",
		"dash/templates/jobs-table.html",
	},
	SrcLogs: []string{
		"dash/templates/page.html",
		"dash/templates/head.html",
		"dash/templates/main-logs.html",
	},
}

func initTemplateFrom(fs embed.FS, filePath string, name string) *template.Template {
	var b []byte
	var err error
	if b, err = fs.ReadFile(filePath); err != nil {
		slog.Error("unable to read from embedded template " + filePath)
	}

	t, err := template.New(name).Parse(string(b))
	if err != nil {
		slog.Error("unable to parse embedded template " + filePath)
	}
	return t
}

// iterate over the embedded source files and init the templates
func initTemplates(fs embed.FS) {
	for i, srcFile := range tplSrcFile {
		tpl[i] = initTemplateFrom(fs, srcFile, i)
	}
	// make a new logs template by parsing the collection of named templates
	dashTpl.logs = template.Must(template.ParseFS(fs, dashTpl.SrcLogs...))
}
