// Package kmsgprobe looks for any smells in the kernel ringbuffer
// See https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg
package kmsgprobe

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/sredog/sre/pkg/analysis"
	gokmsg "github.com/talos-systems/go-kmsg"
)

func CounterToString(counter map[string]int64, sorted bool) string {
	keys := make([]string, 0, len(counter))
	for k := range counter {
		keys = append(keys, k)
	}
	if sorted {
		sort.Strings(keys)
	}
	summary := ""
	sep := ""
	for _, k := range keys {
		val := counter[k]
		summary += fmt.Sprintf("%s%s=%d", sep, k, val)
		sep = ", "
	}
	return summary
}

type KernelRingBufferProbe struct {
	Counter    map[string]int64
	OOMRE      *regexp.Regexp
	OOMVictims map[string]int64
}

const OOMRE = `Killed process (?P<pid>\d+) \((?P<cmd>.+)\) total-vm:(.+), anon-rss:(.+), file-rss:(.+), shmem-rss:(.+)`

// NewKernelRingBufferProbe reads the load average data and returns its representation
func NewKernelRingBufferProbe() (*KernelRingBufferProbe, error) {
	krbp := &KernelRingBufferProbe{
		Counter:    make(map[string]int64),
		OOMRE:      regexp.MustCompile(OOMRE),
		OOMVictims: make(map[string]int64),
	}
	krbp.ReadKernelRingBuffer()
	return krbp, nil
}

func (p *KernelRingBufferProbe) ReadKernelRingBuffer() error {
	reader, err := gokmsg.NewReader()
	if err != nil {
		return err
	}
	defer reader.Close()
	for packet := range reader.Scan(context.TODO()) {
		if packet.Err == nil {
			msg := packet.Message
			// https://github.com/siderolabs/go-kmsg/blob/v0.1.1/message.go#L56
			if msg.Priority > gokmsg.Warning {
				continue
			}
			p.ProcessEvent(msg.Priority.String(), msg.Message)
		}
	}
	return nil
}

func (p *KernelRingBufferProbe) ProcessEvent(priority, message string) {
	val, exists := p.Counter[priority]
	if exists == false {
		p.Counter[priority] = 1
	} else {
		p.Counter[priority] = val + 1
	}
	p.ProcessOOM(message)
}

func (p *KernelRingBufferProbe) ProcessOOM(message string) {
	if strings.Contains(message, "Killed process") {
		matches := p.OOMRE.FindStringSubmatch(message)
		if len(matches) != 7 {
			panic(fmt.Errorf("My regexp game is weak! Wanted 7 finds, got %v in %v", matches, message))
		}
		cmd := matches[2]
		val := p.OOMVictims[cmd]
		p.OOMVictims[cmd] = val + 1
	}
}

const displayFormat = "%v Kernel ring buffer: %s\n"

func (p *KernelRingBufferProbe) Display() string {
	summary := CounterToString(p.Counter, true)
	return fmt.Sprintf(displayFormat,
		emoji.Penguin,
		summary,
	)
}

func (p *KernelRingBufferProbe) Analysis() (observations []*analysis.Observation) {
	observations = append(observations, p.Observations...)
	if len(p.OOMVictims) > 0 {
		var total int64 = 0
		for _, v := range p.OOMVictims {
			total += v
		}
		summary := ""
		if len(p.OOMVictims) < 20 {
			summary = " Counts: " + CounterToString(p.OOMVictims, true)
		}
		observations = append(observations, &analysis.Observation{
			Type:    analysis.Warning,
			Message: fmt.Sprintf("Found %d occurence(s) of OOM killer for %d command(s).%v", total, len(p.OOMVictims), summary),
		})
	}
	observations = append(observations, &analysis.Observation{
		Type:    analysis.Learn,
		Message: "To browse through all kernel ring buffer, use: dmesg --decode --human",
	})
	return
}
