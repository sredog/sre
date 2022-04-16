// Package memory attempts to give useful insight into the system's memory utilisation
// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt
// cat /proc/meminfo
package memory

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
	"github.com/sredog/sre/pkg/format"
)

type MemInfoProvider interface {
	Meminfo() (procfs.Meminfo, error)
}

type MemoryProbe struct {
	Meminfo *procfs.Meminfo
}

// NewMemoryProbe creates an instance of MemoryProbe
func NewMemoryProbe(p MemInfoProvider) (*MemoryProbe, error) {
	mi, err := p.Meminfo()
	if err != nil {
		return nil, err
	}
	la := &MemoryProbe{
		Meminfo: &mi,
	}
	return la, nil
}

const displayFormat = `%v Memory is %s used
Total: %v, available: %v (free: %v, caches: %v, buffers: %v)
Swap total: %v, free: %v
Kernel slab: %v (reclaimable: %v, or %s)
`

func (p *MemoryProbe) Display() string {
	bold := color.New(color.Bold)
	var utilization float64 = 1 - (float64(*p.Meminfo.MemAvailable) / float64(*p.Meminfo.MemTotal))
	textColor := format.ColorForUtilization(utilization)
	var slabReclaimable float64 = (float64(*p.Meminfo.SReclaimable) / float64(*p.Meminfo.Slab))
	var factor uint64 = 1000
	return fmt.Sprintf(displayFormat,
		emoji.ComputerDisk,
		textColor.Sprintf("%0.2f%%", utilization*100),
		// all the values in /proc/meminfo are in kB
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemTotal)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemAvailable)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemFree)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.Cached)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.Buffers)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SwapTotal)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SwapFree)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.Slab)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SReclaimable)),
		textColor.Sprintf("%0.2f%%", slabReclaimable*100),
	)
}

func (p *MemoryProbe) Analysis() (observations []*analysis.Observation) {
	observations = append(observations, &analysis.Observation{
		Type:    analysis.Learn,
		Message: "Have you tried running `cat /proc/meminfo`?",
	})
	return
}
