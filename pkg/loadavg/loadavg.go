// Package loadavg tries to demistify the load averages on Linux
// Kernel implementation: https://github.com/torvalds/linux/blob/master/kernel/sched/loadavg.c
// Helpful blog post: https://www.brendangregg.com/blog/2017-08-08/linux-load-averages.html
package loadavg

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
)

type LoadAvgProvider interface {
	LoadAvg() (*procfs.LoadAvg, error)
}

type LoadAverageProbe struct {
	L *procfs.LoadAvg
}

// NewLoadAverage reads the load average data and returns its representation
func NewLoadAverage(p LoadAvgProvider) (*LoadAverageProbe, error) {
	l, err := p.LoadAvg()
	if err != nil {
		return nil, err
	}
	la := &LoadAverageProbe{
		L: l,
	}
	return la, nil
}

const displayFormat = "%v Load avg: %v (1m), %v (5m), %v (15m)\n"

func (la *LoadAverageProbe) Display() string {
	bold := color.New(color.Bold)
	return fmt.Sprintf(displayFormat,
		emoji.ChartIncreasing,
		bold.Sprintf("%0.2f", la.L.Load1),
		bold.Sprintf("%0.2f", la.L.Load5),
		bold.Sprintf("%0.2f", la.L.Load15),
	)
}

func (la *LoadAverageProbe) Analysis() (observations []*analysis.Observation) {
	epsilon := 0.01
	if la.L.Load1 < epsilon && la.L.Load5 < epsilon && la.L.Load15 < epsilon {
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Note,
			Message: "The system appears idle",
		})
	}
	if la.L.Load1 > la.L.Load5 && la.L.Load1 > la.L.Load15 {
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Note,
			Message: "The load is increasing",
		})
	}
	if la.L.Load1 < la.L.Load5 && la.L.Load1 < la.L.Load15 {
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Note,
			Message: "The load is decreasing",
		})
	}
	observations = append(observations, &analysis.Observation{
		Type:    analysis.Hint,
		Message: "Learn more about load averages https://www.brendangregg.com/blog/2017-08-08/linux-load-averages.html",
	})
	return
}
