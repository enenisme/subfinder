package config

import (
	"flag"
)

type Options struct {
	Domain                     string   // Domain 需要扫描的域名
	DomainsFile                string   // DomainsFile 需要扫描的域名文件
	OutputFile                 string   // OutputFile 最后输出的文件
	TimeoutWithSecond          int      // TimeoutWithSecond 等待源响应的秒数，秒级别
	MaxEnumerateTimeWithMinute int      // MaxEnumerateTimeWithMinute 等待枚举的最长时间，分钟级别
	StdOut                     bool     // StdOut 是否输出到标准输出
	SourcesList                []string // SourcesList 展示的源列表
	IP                         bool     // IP 是否查询IP
	RecordType                 bool     // RecordType 是否查询记录类型
	OutputJson                 bool     // OutputJson 是否输出json格式

	FoundSubdomains []DomainsResult // FoundSubdomains 最后发现的子域名
	FoundNums       int             // FoundNums 最后发现的子域名的数量

}

type DomainsResult struct {
	Domain     string `json:"domain"`      // 子域名
	IP         string `json:"ip"`          // IP 地址
	Source     string `json:"source"`      // 来源
	RecordType string `json:"record_type"` // 记录类型
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Domain, "d", "", "需要扫描的域名")
	flag.StringVar(&options.DomainsFile, "df", "", "需要扫描的域名文件")
	flag.StringVar(&options.OutputFile, "o", "", "输出文件")
	flag.IntVar(&options.TimeoutWithSecond, "t", 10, "等待源响应的秒数，秒级别")
	flag.IntVar(&options.MaxEnumerateTimeWithMinute, "m", 10, "等待枚举的最长时间，分钟级别")

	stdOut := flag.Bool("s", true, "是否输出到标准输出")
	ip := flag.Bool("ip", false, "是否查询IP并输出")
	recordType := flag.Bool("rt", false, "是否查询记录类型并输出")
	outputJson := flag.Bool("oj", false, "是否输出json格式")

	flag.Parse()

	options.StdOut = *stdOut
	options.IP = *ip
	options.RecordType = *recordType
	options.OutputJson = *outputJson

	return options
}

var Logo = `
  _        _ /\ ____        _      __ _           _           
 | |      | |/\/ ___| _   _| |__  / _(_)_ __   __| | ___ _ __ 
 | |   _  | |  \___ \| | | | '_ \| |_| | '_ \ / _ |/  _ \ '__|
 | |__| |_| |   ___) | |_| | |_) |  _| | | | | (_| |  __/ |
 |_____\___/   |____/ \__,_|_.__/|_| |_|_| |_|\__,_|\___|_|
`
