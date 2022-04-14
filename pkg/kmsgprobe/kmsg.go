// Package kmsgprobe looks for any smells in the kernel ringbuffer
// See https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg
package kmsgprobe

import (
	"context"
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/sredog/sre/pkg/analysis"
	gokmsg "github.com/talos-systems/go-kmsg"
)

type KernelRingBufferProbe struct {
}

// NewKernelRingBufferProbe reads the load average data and returns its representation
func NewKernelRingBufferProbe() (*KernelRingBufferProbe, error) {
	reader, err := gokmsg.NewReader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	for packet := range reader.Scan(context.TODO()) {
		if packet.Err == nil {
			msg := packet.Message
			// https://github.com/siderolabs/go-kmsg/blob/v0.1.1/message.go#L56
			if msg.Priority > gokmsg.Warning {
				continue
			}
			// TODO the actual processing of these messages
			fmt.Printf("%v %v", msg.Priority.String(), msg.Message)
		}
	}
	krbp := &KernelRingBufferProbe{}
	return krbp, nil
}

const displayFormat = "%v Kernel ring buffer read\n"

func (la *KernelRingBufferProbe) Display() string {
	return fmt.Sprintf(displayFormat,
		emoji.Penguin,
	)
}

func (la *KernelRingBufferProbe) Analysis() (observations []*analysis.Observation) {
	return
}
