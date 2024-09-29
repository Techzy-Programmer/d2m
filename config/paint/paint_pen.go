package paint

import (
	"github.com/fatih/color"
)

var (
	info    = color.New(color.FgHiCyan)
	err     = color.New(color.FgHiRed)
	warn    = color.New(color.FgHiYellow)
	success = color.New(color.FgHiGreen)
	notice  = color.New(color.FgHiWhite).Add(color.Faint)
)
