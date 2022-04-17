package analysis

// Probe represents probes looking into different aspects of the machine or a process
type Probe interface {
	Displayer
	Analyser
}
