package emails

import (
	"embed"
	"html/template"
)

//go:embed *.tmpl
var tmplFS embed.FS

const PaperAcceptedFile = "paper_accepted.tmpl"
const PaperRejectedFile = "paper_rejected.tmpl"
const PaperPublishedFile = "paper_published.tmpl"

var templates = map[string]*template.Template{
	PaperAcceptedFile:  nil,
	PaperRejectedFile:  nil,
	PaperPublishedFile: nil,
}
