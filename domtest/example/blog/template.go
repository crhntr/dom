package blog

import (
	"embed"
	"text/template"
)

//go:generate muxt generate --receiver-type=App --routes-func=Routes

//go:generate counterfeiter -generate
//counterfeiter:generate -o internal/fake/app.go --fake-name App . RoutesReceiver

//go:embed *.gohtml
var source embed.FS

var templates = template.Must(template.ParseFS(source, "*.gohtml"))
