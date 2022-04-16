// Package memory attempts to give useful insight into the system's memory utilisation
// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt
// cat /proc/meminfo
package memory

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
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

const displayFormat = "%v Memory probe\n"

func (p *MemoryProbe) Display() string {
	return fmt.Sprintf(displayFormat,
		emoji.ComputerDisk,
	)
}

func (p *MemoryProbe) Analysis() (observations []*analysis.Observation) {
	observations = append(observations, &analysis.Observation{
		Type:    analysis.Learn,
		Message: "Have you tried running `cat /proc/meminfo`?",
	})
	return
}
