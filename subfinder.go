package subfinder

import (
	"github.com/enenisme/subfinder/config"
	r "github.com/enenisme/subfinder/pkg/runner"
)

type Subfinder struct {
	Domain     string
	DomainFile string

	Subdomains []config.DomainsResult
}

func NewSubfinder(domain string, domainFile string) *Subfinder {
	return &Subfinder{
		Domain:     domain,
		DomainFile: domainFile,
	}
}

func (s *Subfinder) Run() error {
	opts := config.ParseOptions()
	opts.Domain = s.Domain
	opts.DomainsFile = s.DomainFile
	runner, err := r.NewRunner(opts)
	if err != nil {
		return err
	}

	err = runner.RunEnumeration()
	if err != nil {
		return err
	}

	s.Subdomains = opts.FoundSubdomains
	return nil
}
