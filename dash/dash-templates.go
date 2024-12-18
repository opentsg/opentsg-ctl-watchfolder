package dash

import (
	"embed"
	"html/template"
	"log/slog"
)

type DashTemplateSrc struct {
	main string
}
type DashTemplate struct {
	main *template.Template
}

var tplSrc = DashTemplateSrc{
	main: "dash/templates/main.html",
}
var tpl = DashTemplate{}

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

func initTemplates(fs embed.FS) {
	tpl.main = initTemplateFrom(fs, tplSrc.main, "main")
}
