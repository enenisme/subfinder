package sitedossier

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"time"

	"github.com/enenisme/subfinder/pkg/scrape"
)

var reNext = regexp.MustCompile("<a href=\"([A-Za-z0-9\\/.]+)\"><b>")

type agent struct {
	results chan scrape.Result
	session *scrape.Session
}

type Source struct{}

func (s *Source) Name() string {
	return "sitedossier"
}

func (s *Source) Run(ctx context.Context, domain string, session *scrape.Session) <-chan scrape.Result {
	results := make(chan scrape.Result)

	a := agent{
		session: session,
		results: results,
	}

	go func() {
		_ = a.enumerate(ctx, fmt.Sprintf("http://www.sitedossier.com/parentdomain/%s", domain), s)
	}()
	return results
}

func (a *agent) enumerate(ctx context.Context, baseURL string, s *Source) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := a.session.Get(ctx, baseURL, "", nil)
			if err != nil {
				a.results <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
				return err
			}

			defer func() {
				_ = resp.Body.Close()
				close(a.results)
			}()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				a.results <- scrape.Result{Source: s.Name(), Type: scrape.Error, Error: err}
				return err
			}

			src := string(body)

			for _, match := range a.session.Extractor.FindAllString(src, -1) {
				a.results <- scrape.Result{Source: s.Name(), Type: scrape.Subdomain,
					Value: scrape.DomainInfo{
						Domain:     match,
						SourceType: scrape.API,
					},
				}
			}

			match1 := reNext.FindStringSubmatch(src)
			time.Sleep(time.Duration((3 + rand.Intn(5))) * time.Second)

			if len(match1) > 0 {
				_ = a.enumerate(ctx, "http://www.sitedossier.com"+match1[1], s)
			}
			return nil
		}
	}
}
