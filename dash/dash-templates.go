package dash

import (
	"embed"
	"html/template"
	"log/slog"
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

var tplSrcFile = DashTemplateSrcFileMap{
	"main": "dash/templates/main.html",
	"job":  "dash/templates/job-basic.html",
	"jobs": "dash/templates/jobs-table.html",
}
var tpl = DashTemplateMap{}

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
}
