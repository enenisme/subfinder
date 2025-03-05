package lib

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

type jsonResult struct {
	Input     string `json:"input"`
	Subdomain string `json:"subdomain"`
	Source    string `json:"source"`
}

// WriteDomainToFile 将子域名写入文件
func WriteDomainToFile(domain string, subdomains []string, filename string, json bool) error {
	if json {
		return jsonFile(domain, subdomains, filename)
	}
	return normalFile(subdomains, filename)
}

func normalFile(subdomains []string, filename string) error {
	file, err := createFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	sb := &strings.Builder{}

	for _, subdomain := range subdomains {
		sb.WriteString(subdomain)
		sb.WriteString("\n")

		if _, err = bw.WriteString(sb.String()); err != nil {
			return bw.Flush()
		}
		sb.Reset()
	}

	return bw.Flush()
}

func jsonFile(domain string, subdomains []string, filename string) error {
	file, err := createFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	encode := json.NewEncoder(bw)

	for _, subdomain := range subdomains {
		data := jsonResult{
			Input:     domain,
			Subdomain: subdomain,
			Source:    "subfinder",
		}

		if err = encode.Encode(&data); err != nil {
			return err
		}
	}

	return bw.Flush()
}

// createFile 创建文件
func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_CREATE|os.O_WRONLY, 0644)
}
