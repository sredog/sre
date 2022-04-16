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
	var memoryUtilization float64 = 1 - (float64(*p.Meminfo.MemAvailable) / float64(*p.Meminfo.MemTotal))
	memoryColor := format.ColorForUtilization(memoryUtilization, 0.9, 0.75, 0.5)
	var slabReclaimable float64 = (float64(*p.Meminfo.SReclaimable) / float64(*p.Meminfo.Slab))
	var slabOfTotal float64 = (float64(*p.Meminfo.Slab) / float64(*p.Meminfo.MemTotal))
	slabColor := format.ColorForUtilization(slabOfTotal, 0.2, 0.1, 0.75)
	var factor uint64 = 1000
	return fmt.Sprintf(displayFormat,
		emoji.ComputerDisk,
		memoryColor.Sprintf("%0.2f%%", memoryUtilization*100),
		// all the values in /proc/meminfo are in kB
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemTotal)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemAvailable)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.MemFree)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.Cached)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.Buffers)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SwapTotal)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SwapFree)),
		slabColor.Sprint(humanize.Bytes(factor**p.Meminfo.Slab)),
		bold.Sprint(humanize.Bytes(factor**p.Meminfo.SReclaimable)),
		bold.Sprintf("%0.2f%%", slabReclaimable*100),
	)
}

func (p *MemoryProbe) Analysis() (observations []*analysis.Observation) {
	observations = append(observations, &analysis.Observation{
		Type:    analysis.Learn,
		Message: "Have you tried running `cat /proc/meminfo`?",
	})
	return
}
