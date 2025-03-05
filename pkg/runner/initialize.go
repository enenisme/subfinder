package runner

import "github.com/enenisme/subfinder/pkg/passive"

func (r *Runner) initializePassiveEngine() {
	r.PassiveAgent = passive.NewAgent()
}
