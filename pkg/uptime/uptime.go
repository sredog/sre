// Package uptime provides utilities to read and analyse Linux uptime values
// Learn more at https://man7.org/linux/man-pages/man5/proc.5.html
package uptime

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/sredog/sre/pkg/analysis"
)

const uptimePath = "/proc/uptime"

type UptimeProbe struct {
	Uptime   time.Duration
	Idle     time.Duration
	CPUCount int
}

// NewUptimeProbe reads the uptime data and returns a struct representation
func NewUptimeProbe(CPUCount int) (*UptimeProbe, error) {
	content, err := ioutil.ReadFile(uptimePath)
	if err != nil {
		return nil, err
	}
	// Learn more at https://man7.org/linux/man-pages/man5/proc.5.html
	// /proc/uptime
	// This file contains two numbers (values in seconds): the
	// uptime of the system (including time spent in suspend) and
	// the amount of time spent in the idle process.
	// $ cat /proc/uptime
	// 4932.96 9643.80
	var up, idle float64
	n, err := fmt.Sscanf(string(content), "%f %f", &up, &idle)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		return nil, fmt.Errorf("Expected to read two int values, got %s", content)

	}
	u := &UptimeProbe{
		Uptime:   time.Duration(up*1000) * time.Millisecond,
		Idle:     time.Duration(idle*1000) * time.Millisecond,
		CPUCount: CPUCount,
	}
	return u, nil
}

// Utilization returns ratio spent outside Idle process since boot
// averaged out by the number of cpus
func (u *UptimeProbe) Utilization() float64 {
	return 1 - (float64(u.Idle)/float64(u.CPUCount))/float64(u.Uptime)
}

const displayFormat = `%v Uptime %v
Last boot @ %v
Idle time %v (%v with %d CPUs)
`

func (u *UptimeProbe) Display() string {
	bold := color.New(color.Bold)
	return fmt.Sprintf(displayFormat,
		emoji.AlarmClock,
		bold.Sprint(u.Uptime.String()),
		time.Now().Add(u.Uptime*-1).Format(time.UnixDate),
		// emoji.SleepingFace,
		bold.Sprintf("%0.2f%%", (1-u.Utilization())*100),
		u.Idle,
		u.CPUCount,
	)
}

func (u *UptimeProbe) Analysis() (observations []*analysis.Observation) {
	if u.Uptime < time.Hour*time.Duration(24) {
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Note,
			Message: fmt.Sprintf("This machine restarted recently (%v ago)", u.Uptime),
		})
	}
	return
}
