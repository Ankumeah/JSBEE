package emails

import (
	"fmt"
	"html/template"
)

func InitTemplates() error {
	for file := range templates {
		tmpl, err := template.ParseFS(tmplFS, file)
		if err != nil {
			return fmt.Errorf(
				"Error while parsing file %v: %w",
				file, err,
			)
		}
		templates[file] = tmpl
	}

	return nil
}
