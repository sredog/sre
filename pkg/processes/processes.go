package processes

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
	"github.com/sredog/sre/pkg/format"
)

const PIDMaxPath = "/proc/sys/kernel/pid_max"
const ProcPath = "/proc"

// ReadPIDMax returns the value of pid_max on the system
func ReadPIDMax() (uint64, error) {
	content, err := ioutil.ReadFile(PIDMaxPath)
	if err != nil {
		return 0, err
	}
	// $ cat /proc/sys/kernel/pid_max
	// 4194304
	var limit uint64
	n, err := fmt.Sscanf(string(content), "%d", &limit)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("Expected to read one uint value, got %s", content)
	}
	return limit, nil
}

func CountProcesses() (uint64, error) {
	fd, err := os.Open(ProcPath)
	if err != nil {
		return 0, nil
	}
	defer fd.Close()
	var count uint64
	names, err := fd.Readdirnames(0)
	if err != nil {
		return 0, nil
	}
	for _, name := range names {
		if _, err := strconv.ParseInt(name, 10, 64); err == nil {
			count++
		}
	}
	return count, nil
}

type StatProvider interface {
	Stat() (procfs.Stat, error)
}

type ProcessesProbe struct {
	Stat       *procfs.Stat
	TotalProcs uint64
	PIDMax     uint64
}

// NewProcessesProbe provides insights into the number of processes
func NewProcessesProbe(provider StatProvider) (*ProcessesProbe, error) {
	s, err := provider.Stat()
	if err != nil {
		return nil, err
	}
	limit, err := ReadPIDMax()
	if err != nil {
		return nil, err
	}
	total, err := CountProcesses()
	if err != nil {
		return nil, err
	}
	u := &ProcessesProbe{
		Stat:       &s,
		TotalProcs: total,
		PIDMax:     limit,
	}
	return u, nil
}

const displayFormat = `%v Total processes: %v (%v utilization)
%v running, %v blocked, %v max pid
`

func (p *ProcessesProbe) Display() string {
	bold := color.New(color.Bold)
	var utilization float64 = float64(p.TotalProcs) / float64(p.PIDMax)
	utilisationColor := format.ColorForUtilization(utilization, 0.9, 0.75, 0.5)
	return fmt.Sprintf(displayFormat,
		emoji.RunningShoe,
		bold.Sprintf("%d", p.TotalProcs),
		utilisationColor.Sprintf("%0.2f%%", utilization*100),
		bold.Sprintf("%d", p.Stat.ProcessesRunning),
		bold.Sprintf("%d", p.Stat.ProcessesBlocked),
		bold.Sprintf("%d", p.PIDMax),
	)
}

func (p *ProcessesProbe) Analysis() (observations []*analysis.Observation) {
	var utilization float64 = float64(p.TotalProcs) / float64(p.PIDMax)
	if utilization > 0.75 {
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Warning,
			Message: "You're running out of PIDs - check for zombies",
		})
	}
	return
}
