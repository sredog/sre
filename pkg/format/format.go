package format

import "github.com/fatih/color"

func ColorForUtilization(utilization, high, mid, low float64) (textColor *color.Color) {
	switch {
	case utilization > high:
		textColor = color.New(color.Bold, color.BlinkSlow, color.BgHiRed)
	case utilization > mid:
		textColor = color.New(color.Bold, color.BgHiRed)
	case utilization > low:
		textColor = color.New(color.Bold, color.FgYellow)
	default:
		textColor = color.New(color.Bold)
	}
	return
}
