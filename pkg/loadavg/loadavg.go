// Package loadavg tries to demistify the load averages on Linux
// Kernel implementation: https://github.com/torvalds/linux/blob/master/kernel/sched/loadavg.c
package loadavg

import "github.com/prometheus/procfs"

type LoadAvgProvider interface {
	LoadAvg() (*procfs.LoadAvg, error)
}
