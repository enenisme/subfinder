package resolve

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
)

func QueryIP(domain string) string {
	// 设置 URL
	url := "https://ip.chinaz.com/ipbatch"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("ips", domain)
	_ = writer.WriteField("submore", "查询")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "ip.chinaz.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "multipart/form-data; boundary=--------------------------164773364459923408279548")
	req.Header.Add("Cookie", "qHistory=aHR0cDovL2lwLmNoaW5hei5jb20vaXBiYXRjaC9fSVDmibnph4/mn6Xor6I=")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求解析 IP 地址失败: %v\n", err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("请求解析 IP 读取响应体失败: %v\n", err)
		return ""
	}

	ipRegex := `(?:(?:25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)`
	re := regexp.MustCompile(ipRegex)
	ips := re.FindAllString(string(body), -1)

	if len(ips) == 0 {
		fmt.Println("未找到 IP 地址")
		return ""
	}

	return ips[0]
}
