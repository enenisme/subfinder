package passive

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/enenisme/subfinder/pkg/scrape"
)

// EnumerateSubdomains 枚举给定域的所有子域
func (a *Agent) EnumerateSubdomains(domain string, timeout, maxEnumTimes int) chan scrape.Result {
	result := make(chan scrape.Result)

	go func() {
		defer close(result)
		session, err := scrape.NewSession(domain, timeout)
		if err != nil {
			result <- scrape.Result{Type: scrape.Error, Error: fmt.Errorf("failed to initialize passive session for [%s]: %s", domain, err)}
		}
		defer session.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxEnumTimes)*time.Minute)
		defer cancel()

		var wg sync.WaitGroup

		for sourceName, runner := range a.sources {
			wg.Add(1)

			go func(source string, runner scrape.Source) {
				defer wg.Done()
				for resp := range runner.Run(ctx, domain, session) {
					result <- resp
				}
			}(sourceName, runner)
		}
		wg.Wait()
	}()

	return result
}
