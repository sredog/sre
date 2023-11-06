package cmd

import (
	"github.com/prometheus/procfs"
	"github.com/sredog/sre/pkg/analysis"
	"github.com/sredog/sre/pkg/uptime"
)

type ProbeConfiguration struct {
	ID          string
	Description string
	Aliases     []string
	Build       func(*procfs.FS) (analysis.Probe, error)
}

type ProbeCollectionConfiguration struct {
	ID     string
	Probes []string
}

var probes []*ProbeConfiguration

func init() {
	probes = append(probes, &ProbeConfiguration{
		ID:          "uptime",
		Aliases:     []string{"up"},
		Description: "Look into the systems uptime and idle time",
		Build: func(fs *procfs.FS) (analysis.Probe, error) {
			stat, err := fs.Stat()
			if err != nil {
				panic(err)
			}
			up, err := uptime.NewUptimeProbe(len(stat.CPU))
			if err != nil {
				panic(err)
			}
			return up, nil
		},
	})
}
