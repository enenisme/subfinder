package runner

import "github.com/enenisme/Subfinder/pkg/passive"

func (r *Runner) initializePassiveEngine() {
	r.PassiveAgent = passive.NewAgent()
}
