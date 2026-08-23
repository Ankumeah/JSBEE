package components

import (
	"github.com/a-h/templ"
)

type pageInfo struct {
	title       string
	description string
	content     templ.Component
}
