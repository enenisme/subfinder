package main

import (
	"fmt"
	"time"

	"github.com/enenisme/subfinder/config"
	r "github.com/enenisme/subfinder/pkg/runner"
)

func main() {
	// fmt.Println(config.Logo)
	opts := config.ParseOptions()
	//opts := &config.Options{
	//	Domain:                     "huaun.com",
	//	TimeoutWithSecond:          10,
	//	MaxEnumerateTimeWithMinute: 10,
	//}
	//
	runner, err := r.NewRunner(opts)
	if err != nil {
		fmt.Printf("error creating runner: %v\n", err)
		return
	}

	fmt.Printf("Scanning for subdomains of [%s]: \n", opts.Domain)
	start := time.Now()
	err = runner.RunEnumeration()
	spend := time.Since(start)
	fmt.Printf("found %d subdomains about [%s] in %v\n", opts.FoundNums, opts.Domain, spend)
	if err != nil {
		fmt.Printf("error running enumeration: %v\n", err)
		return
	}

}
