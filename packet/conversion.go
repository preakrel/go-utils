package packet

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"http_fuzzer/common/str"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/dlclark/regexp2"
)

const (
	_              = iota
	Url2Pkt        // url转pack
	PktToGet       // 转get
	PktToPost      // 转post
	PktToMultipart // 转multipart
	PktToJson      //转JSON
	PktToXml       //转XML
	PktToCurl      //数据包转CURL命令
	PktToPython    //数据包转Python代码

	//Url2Packet                    // url转pack
	//Post2Multipart                // post转multipart
	//Multipart2Post                // multipart转post
	//Get2PostNormalFormatPacket    // get转post
	//Get2PostMultipartFormatPacket // get转multipart
	//Post2GetFormatPacket          // post转get

)

var (
	PackFormatError = errors.New("数据包格式不正确")
	PackNilError    = errors.New("数据包是空的")
	URLParseError   = errors.New("url解析失败")
)

type Base struct {
	IsSSL  bool
	Packet string
}

func (b *Base) Byte() []byte {
	return []byte(b.Packet)
}

func (b *Base) judgeMethod(method string) bool {
	s := strings.SplitN(b.Packet, "\n", 2)
	if len(s) == 0 {
		return false
	}
	parts := strings.SplitN(s[0], " ", 2)
	if len(parts) == 0 {
		return false
	}
	return parts[0] == method
}

// Conversion 转化
func Conversion(status uint8, p string, isSSL bool) (Base, error) {
	var err error
	pkt := Base{
		IsSSL:  isSSL,
		Packet: p,
	}
	// 转化操作
	switch status {
	case Url2Pkt:
		pkt, err = URL2Packet(p)
	case PktToGet:
		pkt.Packet, err = PostToGetFormatPacket(p)
	case PktToPost:
		if pkt.judgeMethod(http.MethodPost) {
			var packetLower = strings.ToLower(p)
			if str.StringRegexMatch(packetLower, `^content-type:\s*multipart/`, true) {
				pkt.Packet, err = MultipartToPOST(p)
			} else if str.StringRegexMatch(packetLower, `^content-type:\s*application/xml`, true) { //xml 数据
				pkt.Packet, err = XmlToPOST(p)
			} else if str.StringRegexMatch(packetLower, `^content-type:\s*application/json`, true) { //json 数据
				pkt.Packet, err = JsonToPOST(p)
			}

		} else {
			pkt.Packet, err = GetToPostNormalFormatPacket(p)
		}
	case PktToMultipart:
		if pkt.judgeMethod(http.MethodPost) {
			pkt.Packet, err = PostToMultipart(p)
		} else {
			pkt.Packet, err = GetToPostMultipartFormatPacket(p)
		}
	case PktToJson:
		if pkt.judgeMethod(http.MethodPost) {
			pkt.Packet, err = PostToJson(p)
		} else {
			pkt.Packet, err = GetToPostJsonFormatPacket(p)
		}
	case PktToXml:
		if pkt.judgeMethod(http.MethodPost) {
			pkt.Packet, err = PostToXML(p)
		} else {
			pkt.Packet, err = GetToPostXMLFormatPacket(p)
		}
	case PktToCurl:
		//fmt.Println(p, isSSL)
		translate, err := NewTranslate(p, isSSL)
		if err != nil {
			pkt.Packet = p
			return pkt, err
		}
		pkt.Packet, err = translate.ToCurl()
	case PktToPython:
		translate, err := NewTranslate(p, isSSL)
		if err != nil {
			pkt.Packet = p
			return pkt, err
		}
		pkt.Packet, err = translate.ToPython()
	}

	return pkt, err
}

// URL2Packet URL转数据包
func URL2Packet(link string) (Base, error) {
	var pb Base
	var linkLower = strings.ToLower(link)
	if !strings.HasPrefix(linkLower, "http://") && !strings.HasPrefix(linkLower, "https://") {
		link = "http://" + link
	}
	if strings.HasPrefix(linkLower, "https://") {
		pb.IsSSL = true
	}
	//GET请求
	var packet = ""
	if !strings.Contains(link, `^^^`) {
		uri, err := url.Parse(link)
		if err != nil {
			return pb, err
		}
		port := uri.Port()
		if port == "443" {
			pb.IsSSL = true
		}
		var host = uri.Host //Host是带端口的
		pathAndQuery := uri.RequestURI()
		packet = fmt.Sprintf("GET %s HTTP/1.1\r\n", pathAndQuery)
		packet = packet + fmt.Sprintf("Host: %s\r\n", host)
		packet = fmt.Sprintf("%s%s", packet, "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n")
		packet = fmt.Sprintf("%s%s", packet, "Accept-Language: en\r\n")
		packet = fmt.Sprintf("%s%s", packet, "User-Agent: Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)\r\n")
		packet = fmt.Sprintf("%s%s", packet, "Connection: close\r\n\r\n")
		pb.Packet = packet
	} else { //POST请求,还有multipart
		var actionURL = ""
		var tempStr []string
		typeStr := str.ExtractWithRegex2(linkLower, `(?<=\^\^\^).*?(?=\^\^\^)`, regexp2.RE2)
		if typeStr == "" {
			return pb, fmt.Errorf("url error, url=>%s", link)
		}
		tempStr = strings.Split(link, `^^^`+typeStr+`^^^`)
		actionURL = tempStr[0]
		uri, err := url.Parse(actionURL)
		if err != nil {
			return pb, err
		}
		port := uri.Port()
		if port == "443" {
			pb.IsSSL = true
		}
		var host = uri.Host //Host是带端口的
		pathAndQuery := uri.RequestURI()
		packet = fmt.Sprintf("POST %s HTTP/1.1\r\n", pathAndQuery)
		packet = packet + fmt.Sprintf("Host: %s\r\n", host)
		packet = fmt.Sprintf("%s%s", packet, "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n")
		packet = fmt.Sprintf("%s%s", packet, "Content-Type: application/x-www-form-urlencoded\r\n")
		packet = fmt.Sprintf("%s%s", packet, "Accept-Language: en\r\n")
		packet = fmt.Sprintf("%s%s", packet, "User-Agent: Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)\r\n")
		// 避免url仅包含一个`^^^`时导致panic
		if len(tempStr) >= 2 {
			packet = fmt.Sprintf("%s\r\n%s", packet, tempStr[1])
		} else {
			packet += "\r\n"
		}
		if strings.Contains(linkLower, "^^^multipart^^^") {
			packet, _ = PostToMultipart(packet)
		}
		if strings.Contains(linkLower, "^^^json^^^") {
			packet, _ = PostToJson(packet)
		}
		if strings.Contains(linkLower, "^^^xml^^^") {
			packet, _ = PostToXML(packet)
		}
		pb.Packet = packet
	}
	return pb, nil
}

// PostToMultipart POST转Multipart
func PostToMultipart(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	if str.StringRegexMatch(packet, `^Content-Type:\s*multipart/`, true) {
		return packet, nil
	}

	if str.StringRegexMatch(packet, `^Content-type:\s*application/json`, true) {
		cPacket, err := JsonToPOST(packet)
		if err != nil {
			return packet, nil
		}
		packet = cPacket
	}
	if str.StringRegexMatch(packet, `^Content-type:\s*application/xml`, true) {
		//先转普通的Post再继续
		cPacket, err := XmlToPOST(packet)
		if err != nil {
			return cPacket, nil
		}
		packet = cPacket
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)

	tempArr := strings.SplitN(packet, bodySplitTag, 2)
	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1]
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}
	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), strings.ToLower(`Content-Type:`)) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			}
		}
	}
	boundary := fmt.Sprintf("----%s", str.RandLowerNumberStr(12))
	newPacket = fmt.Sprintf("%sContent-Type: multipart/form-data; boundary=%s%s", newPacket, boundary, fixSplitTag)
	postDataArr := strings.Split(postData, "&")
	var kv []string
	for _, v := range postDataArr {
		v = strings.ReplaceAll(v, "\r", "")
		v = strings.ReplaceAll(v, "\n", "")
		kv = strings.SplitN(v, "=", 2)
		if len(kv) > 1 {
			newPacket = fmt.Sprintf(`%s%s--%s%sContent-Disposition: form-data; name="%s"%s%s%s`, newPacket, fixSplitTag, boundary, fixSplitTag, kv[0], fixSplitTag, fixSplitTag, kv[1])
		} else {
			newPacket = fmt.Sprintf(`%s%s--%s%sContent-Disposition: form-data; name="%s"%s%s`, newPacket, fixSplitTag, boundary, fixSplitTag, kv[0], fixSplitTag, fixSplitTag)
		}
	}
	newPacket = fmt.Sprintf("%s%s--%s--%s", newPacket, fixSplitTag, boundary, fixSplitTag)
	//fmt.Println(newPacket)
	return newPacket, nil
}

// MultipartToPOST Multipart转POST
func MultipartToPOST(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	if !str.StringRegexMatch(packet, `^Content-Type:\s*multipart/`, true) {
		return packet, nil
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)

	tempArr := strings.SplitN(packet, bodySplitTag, 2)
	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1]
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}
	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), `content-type:`) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			} else {
				newPacket = fmt.Sprintf("%sContent-Type: application/x-www-form-urlencoded%s", newPacket, fixSplitTag)
			}
		}
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	r := regexp2.MustCompile(`(?msi)(?<=Content-Disposition:\s*form\-data;\s*name\=\s*("|')).*?(?=---)`, regexp2.Compiled)
	matches, err := r.FindStringMatch(postData)
	if err != nil {
		return packet, PackFormatError
	}
	if matches == nil {
		return packet, errors.New("正则匹配失败")
	}

	var (
		newPostData  = ""
		dataSplitTag string
		key          string
		value        string
	)
	dataSplitTag, _ = getEolSplitStr(matches.String(), `(\r*\n)+`)
	tempArr = strings.Split(matches.String(), dataSplitTag)
	if len(tempArr[0]) > 0 || len(tempArr[1]) > 0 {
		key = tempArr[0]
		value = strings.ReplaceAll(tempArr[1], "\n", "")
		value = strings.TrimRight(value, "\r")
		if strings.HasSuffix(tempArr[0], "'") {
			key = strings.TrimSuffix(tempArr[0], "'")
		}
		if strings.HasSuffix(tempArr[0], `"`) {
			key = strings.TrimSuffix(tempArr[0], `"`)
		}
		newPostData = fmt.Sprintf("%s=%s", key, value)
	}
	for matches != nil {
		matches, err = r.FindNextMatch(matches)
		if err != nil {
			break
		}
		if matches != nil {
			dataSplitTag, _ = getEolSplitStr(matches.String(), `(\r*\n)+`)
			tempArr = strings.Split(matches.String(), dataSplitTag)
			if len(tempArr[0]) > 0 || len(tempArr[1]) > 0 {
				key = tempArr[0]
				value = strings.ReplaceAll(tempArr[1], "\n", "")
				value = strings.TrimRight(value, "\r")
				if strings.HasSuffix(tempArr[0], "'") {
					key = strings.TrimSuffix(tempArr[0], "'")
				}
				if strings.HasSuffix(tempArr[0], `"`) {
					key = strings.TrimSuffix(tempArr[0], `"`)
				}
				if len(newPostData) > 0 {
					newPostData = fmt.Sprintf("%s&%s=%s", newPostData, key, value)
				} else {
					newPostData = fmt.Sprintf("%s=%s", key, value)
				}
			}
		}
	}
	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, newPostData)
	//fmt.Println(newPacket)
	return newPacket, nil
}

func JsonToPOST(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var packetLower = strings.ToLower(packet)

	if !str.StringRegexMatch(packetLower, `^content-type:\s*application/json`, true) {
		return packet, nil
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)

	tempArr := strings.SplitN(packet, bodySplitTag, 2)

	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1] //得到json字符串
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}

	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), `content-type:`) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			} else {
				newPacket = fmt.Sprintf("%sContent-Type: application/x-www-form-urlencoded%s", newPacket, fixSplitTag)
			}
		}
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if !json.Valid([]byte(postData)) { // json format error
		return packet, PackFormatError
	}
	m := make(map[string]interface{})

	d := json.NewDecoder(bytes.NewReader([]byte(postData)))
	d.UseNumber()
	d.Decode(&m)

	newPostData := str.HTTPBuildQuery(m, "", "&", "QUERY_RFC3986")

	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, newPostData)

	return newPacket, nil
}

//PostToJson POST转JSON
func PostToJson(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var packetLower = strings.ToLower(packet)
	if str.StringRegexMatch(packetLower, `^content-type:\s*application/json`, true) {
		return packet, nil
	}
	if str.StringRegexMatch(packetLower, `^content-type:\s*multipart/`, true) {
		//先转普通的Post再继续
		cPacket, err := MultipartToPOST(packet)
		if err != nil {
			return packet, nil
		}
		packet = cPacket
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)

	tempArr := strings.SplitN(packet, bodySplitTag, 2)
	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1]
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}

	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), strings.ToLower(`Content-Type:`)) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			}
		}
	}

	newPacket = fmt.Sprintf("%sContent-Type: application/json;charset=UTF-8%s", newPacket, fixSplitTag)

	var result interface{}
	var err error
	var ok bool
	if str.StringRegexMatch(packetLower, `^content-type:\s*application/xml`, true) {
		//解析XML到MAP
		xmlRet, err := str.NewXmlDecoder(postData).Unmarshal()
		if err != nil {
			return packet, PackFormatError
		}
		result, ok = xmlRet["xml"]
		if !ok {
			result = xmlRet
		}
	} else {
		result, err = str.ParseStr(postData)
	}

	if err != nil {
		return packet, PackFormatError
	}
	r, err := json.Marshal(result)
	if err != nil {
		return packet, PackFormatError
	}
	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, string(r))

	return newPacket, nil
}

//PostToXML POST转XML
func PostToXML(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var packetLower = strings.ToLower(packet)
	if str.StringRegexMatch(packetLower, `^content-type:\s*application/xml`, true) {
		return packet, nil
	}
	if str.StringRegexMatch(packetLower, `^content-type:\s*multipart/`, true) {
		//先转普通的Post再继续
		cPacket, err := MultipartToPOST(packet)
		if err != nil {
			return packet, nil
		}
		packet = cPacket
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)
	tempArr := strings.SplitN(packet, bodySplitTag, 2)
	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1]
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}
	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), strings.ToLower(`Content-Type:`)) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			}
		}
	}

	newPacket = fmt.Sprintf("%sContent-Type: application/xml%s", newPacket, fixSplitTag)
	var result interface{}
	var err error
	if str.StringRegexMatch(packetLower, `content-type:\s*application/json`, true) {
		err = json.Unmarshal([]byte(postData), &result)
		if reflect.TypeOf(result).Kind() == reflect.Map {
			if rXml, ok := result.(map[string]interface{})["xml"]; ok {
				result = rXml
			}
		}
	} else {
		result, err = str.ParseStr(postData)
	}

	if err != nil {
		return packet, PackFormatError
	}
	//Map转Xml
	xmlStr, err := str.NewXmlEncoder(map[string]interface{}{"xml": result}).WithIndent("", "  ").MarshalToString()
	xmlStr = strings.ReplaceAll(xmlStr, "\n", fixSplitTag)
	if err != nil {
		return packet, PackFormatError
	}
	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, xmlStr)
	return newPacket, nil
}

//XmlToPOST XML转POST
func XmlToPOST(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var packetLower = strings.ToLower(packet)

	if !str.StringRegexMatch(packetLower, `^content-type:\s*application/xml`, true) {
		return packet, nil
	}

	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)

	tempArr := strings.SplitN(packet, bodySplitTag, 2)

	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1] //得到xml字符串
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	if len(firstLineArr) == 3 {
		newPacket = fmt.Sprintf("%s %s %s%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		newPacket = fmt.Sprintf("%s %s HTTP/1.1%s", strings.ToUpper(firstLineArr[0]), firstLineArr[1], "\n")
	} else {
		return packet, PackFormatError
	}

	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), `content-type:`) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			} else {
				newPacket = fmt.Sprintf("%sContent-Type: application/x-www-form-urlencoded%s", newPacket, fixSplitTag)
			}
		}
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var peo struct{}
	err := xml.Unmarshal([]byte(postData), &peo)
	if err != nil {
		return packet, PackFormatError
	}
	m, err := str.NewXmlDecoder(postData).Unmarshal()
	if err != nil {
		return packet, PackFormatError
	}
	hbq, ok := m["xml"]
	if !ok {
		hbq = m
	}

	var newPostData string
	if reflect.TypeOf(hbq).Kind() == reflect.Map {
		newPostData = str.HTTPBuildQuery(hbq.(map[string]interface{}), "", "&", "QUERY_RFC3986")
	}

	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, newPostData)

	return newPacket, nil
}

// GetToPostNormalFormatPacket 将GET的数据包转换成POST普通格式的数据包
func GetToPostNormalFormatPacket(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var bodySplitTag = ""
	var fixSplitTag = ""
	var newPacket = ""

	bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)
	tempArr := strings.SplitN(packet, bodySplitTag, 2)
	headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
	if len(tempArr) < 2 {
		return packet, PackFormatError
	}
	postData := tempArr[1]
	headerArr := strings.Split(headerData, "\n")
	firstLineArr := strings.SplitN(headerArr[0], " ", 3)
	var pathAndQueryArr []string
	if len(firstLineArr) == 3 {
		pathAndQueryArr = strings.SplitN(firstLineArr[1], "?", 2)
		newPacket = fmt.Sprintf("POST %s %s%s", pathAndQueryArr[0], firstLineArr[2], "\n")
	} else if len(firstLineArr) == 2 {
		pathAndQueryArr = strings.SplitN(firstLineArr[1], "?", 2)
		newPacket = fmt.Sprintf("POST %s HTTP/1.1%s", pathAndQueryArr[0], "\n")
	} else {
		return packet, PackFormatError
	}
	for _, v := range headerArr[1:] {
		if len(v) > 0 {
			if !strings.HasPrefix(strings.ToLower(v), `content-type:`) {
				newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
			}
		}
	}
	newPacket = fmt.Sprintf("%s%s%s", newPacket, "Content-Type: application/x-www-form-urlencoded", fixSplitTag)
	postData = strings.TrimSpace(postData)
	if len(pathAndQueryArr) == 2 {
		if len(postData) > 0 {
			if strings.HasPrefix(postData, "&") {
				// url参数格式不正确
				postData = strings.TrimLeft(postData, "&")
			}
			if !strings.HasSuffix(postData, "&") {
				postData = fmt.Sprintf("%s&%s", postData, pathAndQueryArr[1])
			} else {
				postData = fmt.Sprintf("%s%s", postData, pathAndQueryArr[1])
			}
		} else {
			postData = pathAndQueryArr[1]
		}
	}
	newPacket = fmt.Sprintf("%s%s%s", newPacket, fixSplitTag, postData)
	return newPacket, nil
}

// GetToPostMultipartFormatPacket 将GET的数据包转换成POST Multipart的数据包
func GetToPostMultipartFormatPacket(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var newPacket = ""
	var err error
	newPacket, err = GetToPostNormalFormatPacket(packet)
	if err != nil {
		return packet, err
	}
	newPacket, err = PostToMultipart(newPacket)
	if err != nil {
		return packet, err
	}
	return newPacket, nil
}

func GetToPostJsonFormatPacket(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var newPacket = ""
	var err error
	//先将GET转成POST
	newPacket, err = GetToPostNormalFormatPacket(packet)
	if err != nil {
		return packet, err
	}
	newPacket, err = PostToJson(newPacket)
	if err != nil {
		return packet, err
	}
	return newPacket, nil
}

func GetToPostXMLFormatPacket(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	var newPacket = ""
	var err error
	//先将GET转成POST
	newPacket, err = GetToPostNormalFormatPacket(packet)
	if err != nil {
		return packet, err
	}
	newPacket, err = PostToXML(newPacket)
	if err != nil {
		return packet, err
	}
	return newPacket, nil
}

// PostToGetFormatPacket 将POST以及POST Multipart的数据包转换为GET的数据包
func PostToGetFormatPacket(packet string) (string, error) {
	if len(packet) == 0 {
		return "", PackNilError
	}
	packetLower := strings.ToLower(packet)
	var newPacket = ""
	if str.StringRegexMatch(packetLower, `^content-type:\s*multipart/`, true) {
		//multipart格式的数据包
		var err error
		newPacket, err = MultipartToPOST(packet)
		if err != nil {
			return packet, nil
		}
		newPacket, err = PostToGetFormatPacket(newPacket)
		if err != nil {
			return packet, nil
		}
	} else if str.StringRegexMatch(packetLower, `^content-type:\s*application/json`, true) { //json 数据
		//JSON格式的数据包

		var err error
		newPacket, err = JsonToPOST(packet)
		if err != nil {
			return packet, nil
		}
		newPacket, err = PostToGetFormatPacket(newPacket)
		if err != nil {
			return packet, nil
		}

	} else if str.StringRegexMatch(packetLower, `^content-type:\s*application/xml`, true) { //xml 数据
		//Xml格式的数据包
		var err error
		newPacket, err = XmlToPOST(packet)
		if err != nil {
			return packet, nil
		}
		newPacket, err = PostToGetFormatPacket(newPacket)
		if err != nil {
			return packet, nil
		}
	} else {
		//普通POST的数据包
		var bodySplitTag = ""
		var fixSplitTag = ""

		bodySplitTag, fixSplitTag = getEolSplitStr(packet, `\r*\n\r*\n`)
		tempArr := strings.SplitN(packet, bodySplitTag, 2)
		headerData := tempArr[0] + fixSplitTag //最后一行没有\r\n或\n了，补上
		if len(tempArr) < 2 {
			return packet, PackFormatError
		}
		postData := strings.ReplaceAll(tempArr[1], "\n", "")
		headerArr := strings.Split(headerData, "\n")
		firstLineArr := strings.SplitN(headerArr[0], " ", 3)
		if len(firstLineArr) == 3 {
			if postData != "" {
				if strings.Contains(firstLineArr[1], "?") {
					newPacket = fmt.Sprintf("GET %s&%s %s%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), fixPostDataToGetParam(postData), firstLineArr[2], "\n")
				} else {
					newPacket = fmt.Sprintf("GET %s?%s %s%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), fixPostDataToGetParam(postData), firstLineArr[2], "\n")
				}
			} else {
				if strings.Contains(firstLineArr[1], "?") {
					newPacket = fmt.Sprintf("GET %s %s%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), firstLineArr[2], "\n")
				} else {
					newPacket = fmt.Sprintf("GET %s %s%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), firstLineArr[2], "\n")
				}
			}

		} else if len(firstLineArr) == 2 {
			if strings.Contains(firstLineArr[1], "?") {
				newPacket = fmt.Sprintf("GET %s&%s HTTP/1.1%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), fixPostDataToGetParam(postData), "\n")
			} else {
				newPacket = fmt.Sprintf("GET %s?%s HTTP/1.1%s", strings.ReplaceAll(firstLineArr[1], " ", "%20"), fixPostDataToGetParam(postData), "\n")
			}
		} else {
			return packet, PackFormatError
		}
		for _, v := range headerArr[1:] {
			if len(v) > 0 {
				if !strings.HasPrefix(strings.ToLower(v), `content-type:`) {
					newPacket = fmt.Sprintf("%s%s%s", newPacket, v, "\n")
				}
			}
		}

		newPacket = fmt.Sprintf("%s%s", newPacket, fixSplitTag)
		return newPacket, nil

	}
	return newPacket, nil
}

func fixPostDataToGetParam(data string) string {
	tmp1 := strings.Split(data, "&")
	var newData string
	for _, v := range tmp1 {
		tmp2 := strings.SplitN(v, "=", 2)
		if len(tmp2) > 1 {
			newData = fmt.Sprintf("%s&%s=%s", newData, escapeParamData(tmp2[0]), escapeParamData(tmp2[1]))
		} else {
			newData = fmt.Sprintf("%s&%s", newData, escapeParamData(tmp2[0]))
		}
	}
	newData = strings.TrimLeft(newData, "&")
	return newData
}

func escapeParamData(data string) string {
	data = strings.ReplaceAll(data, " ", "%20")
	data = strings.ReplaceAll(data, "&", "%26")
	data = strings.ReplaceAll(data, "	", "%09")
	data = strings.ReplaceAll(data, "\r", "%0d")
	data = strings.ReplaceAll(data, "\n", "%0a")
	data = strings.ReplaceAll(data, "+", "%2b")
	data = strings.ReplaceAll(data, "#", "%23")
	data = strings.ReplaceAll(data, "?", "%3f")
	return data
}

func getEolSplitStr(data string, reg string) (string, string) {
	//找Body和Header的分隔符
	var eol string = ""
	if strings.Count(data, "\r\n") >= strings.Count(data, "\n") {
		eol = "\r\n"
	} else {
		eol = "\n"
	}
	regex, err := regexp2.Compile(reg, regexp2.RE2)
	if err != nil {
		return eol, eol
	}
	match, err := regex.FindStringMatch(data)
	if err != nil || match == nil {
		return eol, eol
	}
	return match.String(), eol
}
