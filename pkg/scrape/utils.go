package scrape

import (
	"regexp"
	"sync"
)

var mu sync.Mutex

// NewSubdomainExtractor 创建一个新的子域名提取器
func NewSubdomainExtractor(domain string) (*regexp.Regexp, error) {
	mu.Lock()
	defer mu.Unlock()
	extractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
	if err != nil {
		return nil, err
	}
	return extractor, nil
}
