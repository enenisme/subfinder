package crtsh

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/enenisme/subfinder/pkg/scrape"
)

type Source struct{}

func (s *Source) Name() string {
	return "crtsh"
}

func (s *Source) Run(ctx context.Context, domain string, session *scrape.Session) <-chan scrape.Result {
	result := make(chan scrape.Result)

	go func() {
		defer close(result)
		_ = s.getSubdomainsFromHTTP(ctx, domain, session, result)
	}()

	return result
}

func (s *Source) getSubdomainsFromHTTP(ctx context.Context, domain string, session *scrape.Session, results chan scrape.Result) bool {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := session.Get(ctx, url, "", nil)
	if err != nil {
		results <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
		return false
	}

	src := strings.Replace(string(body), "\\n", " ", -1)

	for _, subdomain := range session.Extractor.FindAllString(src, -1) {
		results <- scrape.Result{Source: s.Name(), Type: scrape.Subdomain,
			Value: scrape.DomainInfo{
				Domain: subdomain,
				//IP:         v.IP,
				SourceType: scrape.CERT,
				//RecordType: v.RecordType,
			}}
	}
	return true
}
