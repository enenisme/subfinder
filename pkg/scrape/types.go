package scrape

import (
	"context"
	"net/http"
	"regexp"
)

type Source interface {
	// Run 将上下文、域名和会话对象作为参数，并返回结果channel
	Run(ctx context.Context, domain string, session *Session) <-chan Result
	// Name 返回源的名称
	Name() string
}

type Session struct {
	HttpClient *http.Client
	Extractor  *regexp.Regexp
}

type Result struct {
	Type   ResultType
	Source string
	Value  DomainInfo
	Error  error
}

type DomainInfo struct {
	Domain     string
	IP         string
	RecordType string
	SourceType SourceType
}

type ResultType int

const (
	Subdomain ResultType = iota
	Error
)

type SourceType string

// api、cert
const (
	API  SourceType = "api"
	CERT SourceType = "cert"
)

var (
	// ErrInitPassiveSession is returned when the session initialization fails
	ErrInitSession = "failed to initialize passive session for scrape"
)
