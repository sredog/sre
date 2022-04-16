package format

import "github.com/fatih/color"

func ColorForUtilization(utilization float64) (textColor *color.Color) {
	switch {
	case utilization > 0.9:
		textColor = color.New(color.Bold, color.BlinkSlow, color.BgHiRed)
	case utilization > 0.8:
		textColor = color.New(color.Bold, color.BgHiRed)
	case utilization > 0.5:
		textColor = color.New(color.Bold, color.FgYellow)
	default:
		textColor = color.New(color.Bold)
	}
	return
}
