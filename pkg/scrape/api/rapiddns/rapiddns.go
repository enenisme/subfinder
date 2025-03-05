package rapiddns

import (
	"context"
	"fmt"
	"io"

	"github.com/enenisme/subfinder/pkg/scrape"
)

type Source struct{}

func (s *Source) Name() string {
	return "rapiddns"
}

func (s *Source) Run(ctx context.Context, domain string, session *scrape.Session) <-chan scrape.Result {
	result := make(chan scrape.Result)

	go func() {
		defer close(result)

		url := fmt.Sprintf("https://rapiddns.io/subdomain/%s", domain)
		resp, err := session.Get(ctx, url, "", nil)
		if err != nil {
			result <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
			return
		}

		src := string(body)
		for _, subdomain := range session.Extractor.FindAllString(src, -1) {
			result <- scrape.Result{Source: s.Name(), Type: scrape.Subdomain,
				Value: scrape.DomainInfo{
					Domain:     subdomain,
					SourceType: scrape.API,
				},
			}

		}
	}()

	return result
}
