package alienvault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/enenisme/Subfinder/pkg/scrape"
)

type Source struct{}

type response struct {
	PassiveDns []struct {
		Host       string `json:"hostname,omitempty" `
		IP         string `json:"address,omitempty"`
		RecordType string `json:"record_type,omitempty"`
	} `json:"passive_dns"`
}

func (s *Source) Name() string {
	return "alienvault"
}

func (s *Source) Run(ctx context.Context, domain string, session *scrape.Session) <-chan scrape.Result {
	result := make(chan scrape.Result)

	go func() {
		defer close(result)

		url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)
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

		var res response
		if err = json.Unmarshal(body, &res); err != nil {
			result <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
			return
		}

		for _, v := range res.PassiveDns {
			result <- scrape.Result{Source: s.Name(), Type: scrape.Subdomain,
				Value: scrape.DomainInfo{
					Domain:     v.Host,
					IP:         v.IP,
					SourceType: scrape.API,
					RecordType: v.RecordType,
				}}
		}

	}()

	return result
}
