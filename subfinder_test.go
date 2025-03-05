package subfinder

import (
	"testing"
)

func TestSubfinder_Run(t *testing.T) {
	subfinder := NewSubfinder("huaun.com", "")
	err := subfinder.Run()
	if err != nil {
		t.Fatalf("error running subfinder: %v", err)
	}
	t.Logf("subdomains: %v", subfinder.Subdomains)
}
