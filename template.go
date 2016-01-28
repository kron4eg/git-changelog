package main

import (
	"text/template"
	"time"
)

var tmpl = `
# Changelog
{{ range . }}
## Version {{ .Version }} @ {{ .Date.Format "2006-01-02" }}
{{ range .Actions }}
### {{ .Action }}
{{ range .Changes }}
  * {{ . }}{{ end }}
{{ end }}{{ end }}`

var (
	changeLogTpl = template.Must(template.New("changelog").Parse(tmpl))
)

type Version struct {
	Version string
	Date    time.Time
	Actions []Action
}

type Action struct {
	Action  string
	Changes []string
}
