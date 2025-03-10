package subfinder

import (
	"testing"
)

func TestSubfinder_Run(t *testing.T) {
	subfinder := NewSubfinder("huaun.com", "", 120)
	err := subfinder.Run()
	if err != nil {
		t.Fatalf("error running subfinder: %v", err)
	}
	t.Logf("subdomains: %v", subfinder.Subdomains)
}
