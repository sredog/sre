package analysis

type ObservationType int64

const (
	Note ObservationType = iota
	Warning
	Issue
	Hint
)

type Observation struct {
	Type  ObservationType
	Lines []string
}

// Analyser is the main interface all probes need to implement
type Analyser interface {
	Analysis() []*Observation
}
