package runner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/enenisme/subfinder/config"
	"github.com/enenisme/subfinder/pkg/lib"
	"github.com/enenisme/subfinder/pkg/resolve"

	"github.com/enenisme/subfinder/pkg/scrape"
)

func (r *Runner) EnumerateDomainForSingle(domain string) error {
	passiveResults := r.PassiveAgent.EnumerateSubdomains(domain, r.Options.TimeoutWithSecond, r.Options.MaxEnumerateTimeWithMinute)

	// 用于去重
	uniqueMap := make(map[string]struct{})

	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		err error
	)

	wg.Add(1)

	go func() {
		defer wg.Done()
		for result := range passiveResults {
			switch result.Type {
			case scrape.Error:
				mu.Lock()
				err = fmt.Errorf("从子域源 [%s] 获取数据错误: %s", result.Source, result.Error)
				mu.Unlock()
				return
			case scrape.Subdomain:
				// 是否是子域
				if !strings.HasSuffix(result.Value.Domain, fmt.Sprintf(".%s", domain)) {
					continue
				}
				subdomain := strings.ReplaceAll(result.Value.Domain, "*.", "")

				// 是否已经存在
				mu.Lock()
				if _, ok := uniqueMap[subdomain]; ok {
					mu.Unlock()
					continue
				}

				uniqueMap[subdomain] = struct{}{}

				wg.Add(1)
				go func(result scrape.Result) {
					defer wg.Done()

					if r.Options.IP {
						if result.Value.IP == "" {
							result.Value.IP = resolve.QueryIP(subdomain)
						}
					}
					if r.Options.RecordType {
						if result.Value.RecordType == "" {
							result.Value.RecordType = "A"
						}
					}

					mu.Lock()
					// 保存结果
					if subdomain != domain {
						r.Options.FoundNums++
						domainResult := config.DomainsResult{
							Domain:     subdomain,
							IP:         result.Value.IP,
							Source:     string(result.Value.SourceType),
							RecordType: result.Value.RecordType,
						}
						if result.Value.IP == "NODOMAIN" {
							domainResult.IP = ""
						}
						r.Options.FoundSubdomains = append(r.Options.FoundSubdomains, domainResult)
					}
					mu.Unlock()
				}(result)

				mu.Unlock()

			}
		}
	}()
	wg.Wait()

	// 如果配置项 OutputFile 不为空，则将结果写入文件
	// TODO: 写入文件，支持是否需要扫描IP
	if r.Options.OutputFile != "" {
		// 写入文件
		domains := make([]string, 0)
		for subdomain := range uniqueMap {
			domains = append(domains, subdomain)
		}
		err = lib.WriteDomainToFile(domain, domains, r.Options.OutputFile, r.Options.OutputJson)
		if err != nil {
			return fmt.Errorf("could not write to file: %v", err)
		}
	}

	if r.Options.StdOut {
		writer := bufio.NewWriter(os.Stdout)
		// 输出到标准输出
		if r.Options.IP {
			if r.Options.RecordType {
				for _, subdomain := range r.Options.FoundSubdomains {
					if subdomain.RecordType == "" {
						subdomain.RecordType = " -"
					}
					_, _ = fmt.Fprintf(writer, "[domain]:%s  [ip]:%s  [record type]:%s\n", subdomain.Domain, subdomain.IP, subdomain.RecordType)
				}
			} else {
				for _, subdomain := range r.Options.FoundSubdomains {
					_, _ = fmt.Fprintf(writer, "[domain]:%s  [ip]:%s \n", subdomain.Domain, subdomain.IP)
				}
			}

		} else {
			for subdomain := range uniqueMap {
				_, _ = fmt.Fprintln(writer, subdomain)
			}
		}
		_ = writer.Flush()
	}

	return err
}

func (r *Runner) EnumerateDomainForMultiple(domains []string) error {
	var wg sync.WaitGroup
	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			if err := r.EnumerateDomainForSingle(domain); err != nil {
				fmt.Printf("error enumerating domain %s: %v\n", domain, err)
				return
			}
		}(domain)
	}
	wg.Wait()

	return nil
}
