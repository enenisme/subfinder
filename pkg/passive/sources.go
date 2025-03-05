package passive

import (
	"github.com/enenisme/subfinder/pkg/scrape"
	"github.com/enenisme/subfinder/pkg/scrape/api/alienvault"
	"github.com/enenisme/subfinder/pkg/scrape/api/crtsh"
	"github.com/enenisme/subfinder/pkg/scrape/api/rapiddns"
)

var DefaultSources = []string{
	// SSL
	"crtsh",

	// API
	"alienvault",
	"sitedossier",
}

// Agent 作为一个被动枚举的代理
type Agent struct {
	sources map[string]scrape.Source
}

// NewAgent 创建一个新的被动枚举代理
func NewAgent() *Agent {
	agent := &Agent{sources: make(map[string]scrape.Source)}
	agent.sources["crtsh"] = &crtsh.Source{}
	agent.sources["alienvault"] = &alienvault.Source{}
	agent.sources["rapiddns"] = &rapiddns.Source{}

	//agent.sources["sitedossier"] = &sitedossier.Source{} // 速度比较慢且十分不稳定，不建议使用

	return agent
}
