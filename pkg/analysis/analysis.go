package analysis

import (
	"fmt"

	"github.com/fatih/color"
)

type ObservationType int64

const (
	Learn ObservationType = iota
	Hint
	Note
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

func (o *Observation) String() string {
	return [...]string{"Learn", "Hint", "Note", "Warning", "Issue"}[o.Type]
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
		textColor = color.New(color.Italic)
	}
	return fmt.Sprintf("%s: %s", textColor.Sprint(o.String()), o.Message)
}
