package analysis

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
)

type ObservationType int64

const (
	Note ObservationType = iota
	Warning
	Issue
	Hint
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
	case Note:
		return "Note"
	case Warning:
		return "Warning"
	case Issue:
		return "Issue"
	}
	return "Note"
}

func (o *Observation) Format() string {
	symbol := emoji.RightArrow
	switch o.Type {
	case Note:
		symbol = emoji.Notebook
	case Warning:
		symbol = emoji.Warning
	case Issue:
		symbol = emoji.Prohibited
	}
	bold := color.New(color.Bold)
	return fmt.Sprintf("%v %s: %s", symbol, bold.Sprint(o.ToString()), o.Message)
}
