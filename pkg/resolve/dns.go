package resolve

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var recordTypes = map[string]uint16{
	"A":     dns.TypeA,
	"AAAA":  dns.TypeAAAA,
	"MX":    dns.TypeMX,
	"CNAME": dns.TypeCNAME,
	"NS":    dns.TypeNS,
	"TXT":   dns.TypeTXT,
}

func QueryDNSConcurrently(domain string) string {
	var wg sync.WaitGroup
	record := make(chan string)
	go func() {
		wg.Wait() // 等待所有查询完成
		close(record)
	}()

	for name, recordType := range recordTypes {
		wg.Add(1)
		go func(name string, t uint16) {
			defer wg.Done()
			if err := queryDNSForType(name, domain, recordType, record); err != nil {
				//fmt.Printf("Error querying %s for type %s: %v\n", domain, name, err)
			}
		}(name, recordType)
	}

	return <-record
}

func queryDNSForType(name, domain string, recordType uint16, record chan string) error {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), recordType)

	c := new(dns.Client)
	c.Timeout = 30 * time.Second  // 增加超时时间到30秒
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return fmt.Errorf("error querying DNS for %s: %v", domain, err)
	}

	for range r.Answer {
		record <- name
	}

	return nil
}
