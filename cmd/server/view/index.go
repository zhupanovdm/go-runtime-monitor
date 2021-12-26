package view

import "html/template"

var Index = template.Must(template.New("index").Parse(index()))

func index() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Title</title>
</head>
<body>
{{range .}}<div><a href="/value/{{.Value.Type}}/{{.Id}}">{{.}}</a></div>{{end}}
</body>
</html>`
}
