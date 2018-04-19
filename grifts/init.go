package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/selfidrone/web_ui/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
