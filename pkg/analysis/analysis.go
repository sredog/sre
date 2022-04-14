package analysis

import (
	"fmt"

	"github.com/fatih/color"
)

type ObservationType int64

const (
	Note ObservationType = iota
	Hint
	Warning
	Issue
)

type Observation struct {
	Type    ObservationType
	Message string
}

// Analyser is the main interface all probes need to implement
type Analyser interface {
	Analysis() []*Observation
}

func (o *Observation) ToString() string {
	switch o.Type {
	default:
		return "Note"
	case Hint:
		return "Hint"
	case Warning:
		return "Warning"
	case Issue:
		return "Issue"
	}
}

func (o *Observation) Format() string {
	var textColor *color.Color
	switch o.Type {
	case Note:
		textColor = color.New(color.Bold)
	case Warning:
		textColor = color.New(color.Bold, color.Underline)
	case Issue:
		textColor = color.New(color.Bold, color.BgHiRed)
	default:
		textColor = color.New(color.BgGreen)
	}
	return fmt.Sprintf("%s: %s", textColor.Sprint(o.ToString()), o.Message)
}
