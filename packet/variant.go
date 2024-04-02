package packet

import (
	"fmt"
	"net/url"
	"strings"
)

type VariantType uint8

const (
	_           VariantType = iota
	URL2Body                // url 转 body
	URL2Cookie              // url 转 cookie
	Body2Cookie             // body 转 cookie
	Body2URL                // body 转 url
	Cookie2URL              // cookie 转 url
	Cookie2Body             // cookie 转 body
)

type variantFactory struct {
	FirstHeaderLine []string // 包头第一行
	OtherHeaderLine []string // 包头其他行
	Body            string   // 包体
	splitTag        string   // 分割符

	originalPacket string // 原始包
}

// 拆包
func (c *variantFactory) resolve(pkt string) error {
	if pkt == "" {
		return PackNilError
	}
	var splitTag string
	if strings.Count(pkt, "\r\n") >= strings.Count(pkt, "\n") {
		splitTag = "\r\n"
	} else {
		splitTag = "\n"
	}
	tempArr := strings.SplitN(pkt, splitTag+splitTag, 2)
	headerData := tempArr[0]
	if len(tempArr) < 2 {
		return PackFormatError
	}
	headerArr := strings.Split(headerData, splitTag)
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) < 2 {
		return PackFormatError
	}

	c.FirstHeaderLine = firstLineArr
	c.OtherHeaderLine = headerArr[1:]
	c.Body = tempArr[1]
	c.splitTag = splitTag
	c.originalPacket = pkt
	return nil
}

// 打包
func (c *variantFactory) recombine() string {
	//firstLine := strings.Join(c.FirstHeaderLine, " ")
	//otherLines := strings.Join(c.OtherHeaderLine, c.splitTag)
	var otherLines, firstLine string
	for _, s := range c.FirstHeaderLine {
		if s == "" {
			continue
		}
		if firstLine == "" {
			firstLine = fmt.Sprintf("%s", s)
		} else {
			firstLine = fmt.Sprintf("%s %s", firstLine, s)
		}
	}

	for _, s := range c.OtherHeaderLine {
		if s == "" {
			continue
		}
		if otherLines == "" {
			otherLines = fmt.Sprintf("%s", s)
		} else {
			otherLines = fmt.Sprintf("%s%s%s", otherLines, c.splitTag, s)
		}

	}
	newPacket := fmt.Sprintf("%s%s%s%s%s", firstLine, c.splitTag, otherLines, c.splitTag+c.splitTag, c.Body)
	return newPacket
}

// 设置url
func (c *variantFactory) setURL(value string) {
	c.FirstHeaderLine[1] = value
}

// 获取url
func (c *variantFactory) getURL() string {
	return c.FirstHeaderLine[1]
}

// 设置cookies
func (c *variantFactory) setCookie(value string) {
	if strings.Contains(strings.ToLower(c.originalPacket), `cookie:`) {
		// 存在cookie字段
		for i, v := range c.OtherHeaderLine {
			if len(v) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(v), `cookie:`) {

				part := strings.SplitN(v, ":", 2)
				if len(part) > 1 {
					part[1] = fmt.Sprintf("%s%s%s", part[1], "; ", value)
					v = strings.Join(part, ":")
					c.OtherHeaderLine[i] = v
				}
			}

		}
	} else {
		// 不存在cookie字段
		cookieHeader := fmt.Sprintf(`Cookie: %s`, value)
		c.OtherHeaderLine = append(c.OtherHeaderLine, cookieHeader)
	}
}

// 提取cookies
func (c *variantFactory) extractCookie() string {
	if strings.Contains(strings.ToLower(c.originalPacket), `cookie:`) {
		for i, h := range c.OtherHeaderLine {
			if strings.HasPrefix(strings.ToLower(h), `cookie:`) {
				part := strings.SplitN(h, ":", 2)
				if len(part) > 1 {
					c.OtherHeaderLine[i] = ""
					return part[1]
				}
				return ""
			}
		}
	}

	return ""
}

// Variant 变形
func Variant(status uint8, p string) (string, error) {
	var err error
	var pkt string

	// 转化操作
	switch VariantType(status) {
	case URL2Body:
		pkt, err = URLToBody(p)
	case URL2Cookie:
		pkt, err = URLToCookie(p)
	case Body2Cookie:
		pkt, err = BodyToCookie(p)
	case Body2URL:
		pkt, err = BodyToURL(p)
	case Cookie2URL:
		pkt, err = CookieToURL(p)
	case Cookie2Body:
		pkt, err = CookieToBody(p)
	}

	return pkt, err
}

// URLToBody url 转 body
func URLToBody(pkt string) (string, error) {
	return GetToPostNormalFormatPacket(pkt)
}

// URLToCookie url 转 cookie
func URLToCookie(pkt string) (string, error) {
	c := variantFactory{}
	err := c.resolve(pkt)
	if err != nil {
		return pkt, err
	}

	urlParse, err := url.Parse(c.getURL())
	if err != nil {
		return pkt, URLParseError
	}
	rawQuery := urlParse.RawQuery
	urlParse.RawQuery = ""
	c.setURL(urlParse.String())

	// 转化成cookie参数
	params := strings.Split(rawQuery, "&")
	cookieValue := strings.Join(params, "; ")
	c.setCookie(cookieValue)
	newPacket := c.recombine()

	return newPacket, nil
}

// BodyToCookie body 转 cookie
func BodyToCookie(pkt string) (string, error) {
	c := variantFactory{}
	err := c.resolve(pkt)
	if err != nil {
		return pkt, err
	}

	if c.Body == "" {
		return pkt, nil
	}

	params := strings.Split(c.Body, "&")
	cookieValue := strings.Join(params, "; ")
	c.setCookie(cookieValue)
	c.Body = ""
	newPacket := c.recombine()

	return newPacket, nil
}

// BodyToURL body 转 url
func BodyToURL(pkt string) (string, error) {
	c := variantFactory{}
	err := c.resolve(pkt)
	if err != nil {
		return pkt, err
	}

	if c.Body == "" {
		return pkt, nil
	}

	urlParse, err := url.Parse(c.getURL())
	if err != nil {
		return pkt, URLParseError
	}
	paramsURL := urlParse.Query()
	paramsBody := strings.Split(c.Body, "&")
	for _, param := range paramsBody {
		p := strings.TrimSpace(param)
		kv := strings.SplitN(p, "=", 2)
		paramsURL.Set(kv[0], kv[1])
	}
	urlParse.RawQuery = paramsURL.Encode()
	c.setURL(urlParse.String())
	c.Body = ""
	newPacket := c.recombine()

	return newPacket, nil
}

// CookieToURL cookie 转 url
func CookieToURL(pkt string) (string, error) {
	c := variantFactory{}
	err := c.resolve(pkt)
	if err != nil {
		return pkt, err
	}

	cookiesParams := strings.Split(c.extractCookie(), "; ")
	urlParse, err := url.Parse(c.getURL())
	if err != nil {
		return pkt, URLParseError
	}
	paramsURL := urlParse.Query()
	for _, param := range cookiesParams {
		p := strings.TrimSpace(param)
		kv := strings.SplitN(p, "=", 2)
		paramsURL.Set(kv[0], kv[1])
	}
	urlParse.RawQuery = paramsURL.Encode()
	c.setURL(urlParse.String())
	newPacket := c.recombine()

	return newPacket, nil
}

// CookieToBody cookie 转 body
func CookieToBody(pkt string) (string, error) {
	c := variantFactory{}
	err := c.resolve(pkt)
	if err != nil {
		return pkt, err
	}

	cookiesParams := strings.Split(c.extractCookie(), "; ")
	for _, param := range cookiesParams {
		p := strings.TrimSpace(param)
		if c.Body == "" {
			c.Body = fmt.Sprintf("%s", p)
		} else {
			c.Body = fmt.Sprintf("%s&%s", c.Body, p)
		}

	}
	newPacket := c.recombine()

	return newPacket, nil
}
