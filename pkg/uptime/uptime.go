// Package uptime provides utilities to read and analyse Linux uptime values
// Learn more at https://man7.org/linux/man-pages/man5/proc.5.html
package uptime

import (
	"fmt"
	"io/ioutil"
	"time"
)

const uptimePath = "/proc/uptime"

type Uptime struct {
	Uptime   time.Duration
	Idle     time.Duration
	CPUCount int
}

// NewUptime reads the uptime data and returns a struct representation
func NewUptime(CPUCount int) (*Uptime, error) {
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
	u := &Uptime{
		Uptime:   time.Duration(up*1000) * time.Millisecond,
		Idle:     time.Duration(idle*1000) * time.Millisecond,
		CPUCount: CPUCount,
	}
	return u, nil
}
