package cpu

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
	"github.com/sredog/sre/pkg/format"
)

// CPUTotalTime adds up all categories
func CPUTotalTime(cpu *procfs.CPUStat) float64 {
	return cpu.User +
		cpu.Nice +
		cpu.System +
		cpu.Idle +
		cpu.Iowait +
		cpu.IRQ +
		cpu.SoftIRQ +
		cpu.Steal +
		cpu.Guest +
		cpu.GuestNice
}

type StatProvider interface {
	Stat() (procfs.Stat, error)
}

type CPUProbe struct {
	Stat *procfs.Stat
}

// NewCPUProbe provides insights into CPU utilization
func NewCPUProbe(provider StatProvider) (*CPUProbe, error) {
	s, err := provider.Stat()
	if err != nil {
		return nil, err
	}
	u := &CPUProbe{
		Stat: &s,
	}
	return u, nil
}

const displayFormat = `%v %v CPUs at %v utilization
User: %v (niced %v), system: %v, stolen: %v, idle: %v
`

func (p *CPUProbe) Display() string {
	bold := color.New(color.Bold)
	total := CPUTotalTime(&p.Stat.CPUTotal)
	var utilization float64 = 1 - (float64(p.Stat.CPUTotal.Idle) / total)
	utilisationColor := format.ColorForUtilization(utilization, 0.95, 0.85, 0.5)
	return fmt.Sprintf(displayFormat,
		emoji.Fire,
		bold.Sprintf("%d", len(p.Stat.CPU)),
		utilisationColor.Sprintf("%0.2f%%", utilization*100),
		bold.Sprintf("%0.2f%%", p.Stat.CPUTotal.User/total*100),
		bold.Sprintf("%0.2f%%", p.Stat.CPUTotal.Nice/total*100),
		bold.Sprintf("%0.2f%%", p.Stat.CPUTotal.System/total*100),
		bold.Sprintf("%0.2f%%", p.Stat.CPUTotal.Steal/total*100),
		bold.Sprintf("%0.2f%%", p.Stat.CPUTotal.Idle/total*100),
	)
}

func (p *CPUProbe) Analysis() (observations []*analysis.Observation) {
	// TODO detect when CPUs are not equally busy
	return
}
