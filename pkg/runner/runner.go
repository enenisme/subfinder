package runner

import (
	"github.com/enenisme/subfinder/config"
	"github.com/enenisme/subfinder/pkg/lib"
	"github.com/enenisme/subfinder/pkg/passive"
)

type Runner struct {
	Options      *config.Options
	PassiveAgent *passive.Agent
}

func NewRunner(options *config.Options) (*Runner, error) {
	opts := options
	opts.SourcesList = passive.DefaultSources

	runner := &Runner{Options: opts}
	runner.initializePassiveEngine()
	return runner, nil
}

func (r *Runner) RunEnumeration() error {
	if r.Options.DomainsFile != "" {
		subdomains, err := lib.ReadDomainsFromFile(r.Options.DomainsFile)
		if err != nil {
			return err
		}
		return r.EnumerateDomainForMultiple(subdomains)
	}
	return r.EnumerateDomainForSingle(r.Options.Domain)
}
