package kmsgprobe

import "testing"

func TestProcessOOMBasic(t *testing.T) {
	p, _ := NewKernelRingBufferProbe()
	p.ProcessEvent("error", "Killed process 102503 (chrome) total-vm:911152kB, anon-rss:104748kB, file-rss:0kB, shmem-rss:1968kB")
	if len(p.OOMVictims) != 1 {
		t.Errorf("Expected 1 OOM event to be detected, got %v", p.OOMVictims)
	}
	val, contains := p.OOMVictims["chrome"]
	if contains != true {
		t.Errorf("Expected to detect OOM process")
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
}
