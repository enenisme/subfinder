package config

import (
	"flag"
	"os"
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

	// 创建新的 FlagSet 来避免与其他包的 flags 冲突
	flags := flag.NewFlagSet("subfinder", flag.ContinueOnError)

	flags.StringVar(&options.Domain, "d", "", "需要扫描的域名")
	flags.StringVar(&options.DomainsFile, "df", "", "需要扫描的域名文件")
	flags.StringVar(&options.OutputFile, "o", "", "输出文件")
	flags.IntVar(&options.TimeoutWithSecond, "t", 10, "等待源响应的秒数，秒级别")
	flags.IntVar(&options.MaxEnumerateTimeWithMinute, "m", 10, "等待枚举的最长时间，分钟级别")

	stdOut := flags.Bool("s", true, "是否输出到标准输出")
	ip := flags.Bool("ip", false, "是否查询IP并输出")
	recordType := flags.Bool("rt", false, "是否查询记录类型并输出")
	outputJson := flags.Bool("oj", false, "是否输出json格式")

	// 只解析未解析的参数
	if !flags.Parsed() {
		_ = flags.Parse(os.Args[1:])
	}

	options.StdOut = *stdOut
	options.IP = *ip
	options.RecordType = *recordType
	options.OutputJson = *outputJson

	return options
}
