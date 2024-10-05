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

// 0: info, 1: err, 2: warn, 3: success, 4: notice
func GetLogLevel(logger *color.Color) uint {
	var logLvl uint = 0

	switch logger {
	case err:
		logLvl = 1
	case warn:
		logLvl = 2
	case success:
		logLvl = 3
	case notice:
		logLvl = 4
	}

	return logLvl
}
