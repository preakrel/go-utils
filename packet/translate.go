package packet

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func NewTranslate(msg string, isSSL bool) (*Translater, error) {
	if len(msg) == 0 {
		return nil, PackNilError
	}
	c := &Translater{
		ssl:       isSSL,
		msg:       msg,
		msgReader: bufio.NewReader(strings.NewReader(msg)),
	}
	err := c.parse()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type Translater struct {
	ssl        bool
	msg        string
	msgReader  *bufio.Reader
	request    *http.Request
	builder    strings.Builder
	requestURL string //请求地址
	storeBody  string
}

func (c *Translater) parse() error {
	var err error
	c.request, err = http.ReadRequest(c.msgReader)
	if err != nil {
		return errors.New("http message is illegal")
	}
	err = c.fix()
	return err
}

func (c *Translater) check() error {
	if c == nil || c.request == nil {
		return errors.New("packet is required")
	}
	if c.request.Host == "" {
		return errors.New("host is required")
	}
	if c.request.Method == "" {
		return errors.New("method is required")
	}
	schema := c.request.URL.Scheme
	requestURL := c.request.URL.String()
	if schema == "" {
		if c.ssl {
			schema = "https"
		} else {
			schema = "http"
		}
		tmp := strings.SplitN(c.request.Host, ":", 2)
		var (
			host string
			port string
		)
		if len(tmp) == 2 {
			host = tmp[0]
			port = tmp[1]
		} else {
			host = tmp[0]
			port = ""
		}
		if port == "443" {
			requestURL = "https://" + host + c.request.RequestURI
		} else if port == "80" {
			requestURL = "http://" + host + c.request.RequestURI
		} else {
			requestURL = schema + "://" + c.request.Host + c.request.RequestURI
		}
	}
	c.requestURL = requestURL
	return nil
}

func (c *Translater) fix() error {
	var realLen int64 = 0
	var splitTag = ""
	if strings.Count(c.msg, "\r\n") >= strings.Count(c.msg, "\n") {
		splitTag = "\r\n"
	} else {
		splitTag = "\n"
	}
	tempArr := strings.SplitN(c.msg, splitTag+splitTag, 2)
	if len(tempArr) > 1 {
		c.storeBody = tempArr[1]
		realLen = int64(len(tempArr[1]))
	}
	if c.request.ContentLength == realLen {
		return nil
	}
	compile := regexp.MustCompile(`Content-Length: \d*`)
	findString := compile.FindString(c.msg)
	if findString != "" {
		c.msg = compile.ReplaceAllString(c.msg, fmt.Sprintf("Content-Length: %d", realLen))
	} else if realLen > 0 { //自动填充Content-Length
		c.msg = fmt.Sprintf("%s%s%s", fmt.Sprintf("%s%sContent-Length: %s%s", tempArr[0], splitTag, strconv.FormatInt(realLen, 10), splitTag), splitTag, tempArr[1])
	}
	c.msgReader = bufio.NewReader(strings.NewReader(c.msg))
	return c.parse()
}

// doCurl 构建命令行
func (c *Translater) ToCurl() (string, error) {
	var err error
	err = c.check()
	if err != nil {
		return "", err
	}
	c.builder.Reset()
	c.commandMajor()
	c.commandHeaders()

	err = c.commandBody()
	if err != nil {
		// c.builder.WriteString(fmt.Sprintf(`--data-raw %s`, fmt.Sprintf("%#v", c.storeBody)))
		return "", err
	}
	return c.builder.String(), nil
}

// doPython 构建Python代码
func (c *Translater) ToPython() (string, error) {
	var err error
	err = c.check()
	if err != nil {
		return "", err
	}
	c.builder.Reset()
	c.pythonBeg()
	c.pythonUrl()
	c.pythonHeaders()
	c.pythonCookies()
	err = c.pythonBody()
	if err != nil {
		return "", err
	}
	c.pythonReq()
	return c.builder.String(), nil
}

func (c *Translater) pythonBeg() {
	c.builder.WriteString("#!/usr/bin/python")
	c.builder.WriteString("\n")
	c.builder.WriteString("# -*- coding: UTF-8 -*-")
	c.builder.WriteString("\n")
	c.builder.WriteString("import requests")
	c.builder.WriteString("\n")
	if strings.Contains(c.request.Header.Get("Content-Type"), "/json") {
		c.builder.WriteString("import json")
	}
	c.builder.WriteString("\n\n")
}

func (c *Translater) pythonUrl() {
	c.builder.WriteString(fmt.Sprintf(`url = %s`, quote(c.requestURL)))
	c.builder.WriteString("\n")

}

func (c *Translater) pythonHeaders() {
	headerMap := make(map[string]string)
	for name, vals := range c.request.Header {
		if reflect.DeepEqual(name, "Cookie") {
			continue
		}
		for _, val := range vals {
			headerMap[name] = val
		}
	}
	headerBytes, _ := json.Marshal(headerMap)
	c.builder.WriteString(fmt.Sprintf(`headers = %s`, string(headerBytes)))
	c.builder.WriteString("\n")
}

func (c *Translater) pythonCookies() {
	cookiesMap := make(map[string]string)
	for _, val := range c.request.Cookies() {
		cookiesMap[val.Name] = val.Value

	}
	cookiesBytes, _ := json.Marshal(cookiesMap)
	c.builder.WriteString(fmt.Sprintf(`cookies = %s`, string(cookiesBytes)))
	c.builder.WriteString("\n")
}

func (c *Translater) pythonBody() error {
	all, err := io.ReadAll(c.request.Body)
	if err != nil {
		return fmt.Errorf("parse body failed,%s", err)
	}

	c.builder.WriteString(fmt.Sprintf(`data = %s`, quote(string(all))))
	c.builder.WriteString("\n")
	if strings.Contains(c.request.Header.Get("Content-Type"), "/json") {
		if !json.Valid(all) {
			return fmt.Errorf("parse json failed,%s", err)
		}
		c.builder.WriteString(`jsonData = json.loads(data)`)
	}
	c.builder.WriteString("\n\n")
	return nil
}

// func (c *Converter) pythonForm() error {
// 	err := c.request.ParseMultipartForm(c.request.ContentLength)
// 	if err != nil {
// 		return fmt.Errorf("parse multipart/form-data failed,%s", err)
// 	}
// 	ms := make(map[string]interface{})
// 	if c.request.MultipartForm.Value != nil {
// 		for name, vals := range c.request.MultipartForm.Value {
// 			for _, val := range vals {
// 				ms[name] = recast(val, true)
// 			}
// 		}
// 	}
// 	msBytes, _ := json.Marshal(ms)
// 	c.builder.WriteString(fmt.Sprintf(`data = %s`, string(msBytes)))
// 	c.builder.WriteString("\n\n")
// 	return nil
// }
// func (c *Converter) pythonFormurl() error {
// 	err := c.request.ParseForm()
// 	if err != nil {
// 		return fmt.Errorf("parse form-urlencoded failed,%s", err)
// 	}
// 	ms := make(map[string]interface{})
// 	for name, vals := range c.request.PostForm {
// 		for _, val := range vals {
// 			ms[name] = recast(val, true)
// 		}
// 	}
// 	msBytes, _ := json.Marshal(ms)
// 	c.builder.WriteString(fmt.Sprintf(`data = %s`, string(msBytes)))
// 	c.builder.WriteString("\n\n")
// 	return nil
// }

func (c *Translater) pythonReq() error {
	var html string
	if strings.Contains(c.request.Header.Get("Content-Type"), "/json") {
		html = ", json=jsonData"
	}
	c.builder.WriteString(fmt.Sprintf(`html = requests.request("%s" ,url, headers=headers, verify=False, cookies=cookies , data=data%s)`, strings.ToUpper(c.request.Method), html))
	c.builder.WriteString("\n\n")
	c.builder.WriteString("print(len(html.text))\n")
	c.builder.WriteString("print(html.text)\n")
	return nil
}

// command
func (c *Translater) commandMajor() {
	c.builder.WriteString(fmt.Sprintf(`curl -i -s -k -X $"%s" \`, c.request.Method))
	c.builder.WriteString("\n")

}

func (c *Translater) commandHeaders() {
	for name, vals := range c.request.Header {
		for _, val := range vals {
			c.builder.WriteString(fmt.Sprintf(`-H $%s \`, quote(fmt.Sprintf("%s: %s", name, val))))
			c.builder.WriteString("\n")
		}
	}
}

func (c *Translater) commandBody() error {
	all, err := io.ReadAll(c.request.Body)
	if err != nil {
		return fmt.Errorf("parse body failed,%s", err)
	}
	if len(all) > 0 {
		c.builder.WriteString(fmt.Sprintf(`--data-binary  $%s \`, quote(string(all))))
		c.builder.WriteString("\n")
	}
	c.builder.WriteString(fmt.Sprintf(`$%s`, quote(c.requestURL)))
	c.builder.WriteString("\n")

	return nil
}

func quote(s string) string {
	str := strconv.QuoteToASCII(s)
	str = strings.Replace(str, "\\n", "\\x0a", -1)
	str = strings.Replace(str, "\\r", "\\x0d", -1)
	// str = strings.Replace(str, " ", "\\x20", -1)
	return str
}
