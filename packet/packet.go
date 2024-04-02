package packet

import (
	"encoding/base64"
	"fmt"
	"http_fuzzer/common/str"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

type HttpPacket struct {
	OriginalPacket       string //原始的数据包
	MarkedPacket         string //标记的测试包
	CurrentPayLoadPacket string //当前的变体包

	Valid bool //该包是否有效

	Scheme  string //http或者https
	Method  string //请求方式
	HttpVer string //HTTP协议版本

	DestIP   string //发往的目标IP
	DestPort int    //发往的目标端口

	URL          string //完整URL
	Host         string //域名/主机,注意不带端口
	HostAndPort  string //主机+端口
	Port         int    //端口
	PathAndQuery string
	Path         string
	Query        map[string][]string
	RawQuery     string
	Fragment     string
	AuthUser     string //auth 认证用户名
	AuthPwd      string //auth 认证密码

	UserAgent string

	Cookies    map[string][]string
	RawCookies string

	Headers    map[string][]string //不含UA、Cookie头部
	RawHeaders string

	DupHeaders    map[string][]string //多余的头部
	RawDupHeaders string

	PostData      string //请求body数据
	PostDataBytes []byte //请求body的bytes形式

	BodySplit string //Header和Body的分隔符

	IsSSL bool //是否
}

// MapArrayGet 从Map数组中查询一个值
func MapArrayGet(v map[string][]string, key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// MapArraySet 从Map数组中设置一个值
func MapArraySet(v map[string][]string, key, value string) {
	v[key] = []string{value}
}

// MapArrayAdd 从Map数组中添加一个值
func MapArrayAdd(v map[string][]string, key, value string) {
	v[key] = append(v[key], value)
}

// MapArrayDel 从Map数组中删除一个值
func MapArrayDel(v map[string][]string, key string) {
	delete(v, key)
}

// MapArrayToString 将Map数组中的数据拼接成字符串
func MapArrayToString(v map[string][]string, colStr, lineStr string) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		//keyEscaped := url.QueryEscape(k)
		keyEscaped := k
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString(colStr)
			}
			buf.WriteString(keyEscaped)
			if v == "§NIL§" {
				continue
			}
			buf.WriteString(lineStr)
			//buf.WriteString(url.QueryEscape(v))
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// MapArrayEncodeToString 将Map数组中的数据拼接成字符串，转义k和v
func MapArrayEncodeToString(v map[string][]string, colStr, lineStr string) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString(colStr)
			}
			buf.WriteString(keyEscaped)
			buf.WriteString(lineStr)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

func (pkt *HttpPacket) SetAuthUser(user string) {
	pkt.AuthUser = user
}

func (pkt *HttpPacket) SetAuthPwd(pwd string) {
	pkt.AuthPwd = pwd
}

// SetSSL 设置SSL
func (pkt *HttpPacket) SetSSL(value bool) {
	pkt.IsSSL = value
	//同时修正URL
	if value {
		pkt.Scheme = "https"
		pkt.URL = str.RegexReplace(pkt.URL, "^http://", "https://", regexp2.IgnoreCase)
	} else {
		pkt.Scheme = "http"
		pkt.URL = str.RegexReplace(pkt.URL, "^https://", "http://", regexp2.IgnoreCase)
	}
}

// GetSSL 获取SSL
func (pkt *HttpPacket) GetSSL() bool {
	return pkt.IsSSL
}

// SetMethod 设置Method
func (pkt *HttpPacket) SetMethod(value string) {
	pkt.Method = value
}

// GetMethod 获取Method
func (pkt *HttpPacket) GetMethod() string {
	return pkt.Method
}

// SetHttpVer 设置Http协议版本
func (pkt *HttpPacket) SetHttpVer(value string) {
	pkt.HttpVer = value
}

// GetHttpVer 获取Http协议版本
func (pkt *HttpPacket) GetHttpVer() string {
	return pkt.HttpVer
}

// SetDestIP 设置要发往目标主机的IP
func (pkt *HttpPacket) SetDestIP(value string) {
	pkt.DestIP = value
}

// GetDestIP 获取要发往目标主机的IP
func (pkt *HttpPacket) GetDestIP() string {
	return pkt.DestIP
}

// SetDestPort 设置要发往目标主机的端口
func (pkt *HttpPacket) SetDestPort(value int) {
	pkt.DestPort = value
}

// GetDestPort 获取要发往目标主机的端口
func (pkt *HttpPacket) GetDestPort() int {
	return pkt.DestPort
}

// SetURL 设置请求的URL
func (pkt *HttpPacket) SetURL(value string) {
	uri, err := url.Parse(value)
	//同时修正其他的相关参数
	if err == nil {
		pkt.URL = value
		pkt.Scheme = uri.Scheme
		if uri.Scheme == "http" {
			pkt.IsSSL = false
		} else {
			pkt.IsSSL = true
		}
		if strings.Contains(uri.Host, ":") {
			tmp := strings.SplitN(uri.Host, ":", 2)
			pkt.Host = tmp[0]
			pkt.Port = str.String2IntWithDefaultValue(tmp[1], 80)
			pkt.HostAndPort = uri.Host
		} else {
			pkt.Host = uri.Host
			if pkt.IsSSL {
				pkt.Port = 443
				pkt.HostAndPort = uri.Host
			} else {
				pkt.Port = 80
				pkt.HostAndPort = uri.Host
			}
		}
		pkt.Path = uri.Path //uri.RawPath 有问题
		if len(pkt.Path) == 0 {
			pkt.Path = "/"
		}
		if len(uri.RawQuery) > 0 {
			pkt.PathAndQuery = pkt.Path + "?" + uri.RawQuery
		} else {
			pkt.PathAndQuery = pkt.Path
		}
		pkt.Query = uri.Query()
		pkt.RawQuery = uri.RawQuery
		pkt.Fragment = uri.RawFragment
	}
}

// GetURL 获取请求的URL
func (pkt *HttpPacket) GetURL() string {
	return pkt.URL
}

// SetHost 设置请求的Host,不带端口
func (pkt *HttpPacket) SetHost(value string) {
	pkt.Host = value
	// 同时修正其他的相关参数
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
	pkt.SetURL(pkt.URL)
}

// GetHost 获取请求的Host,不带端口
func (pkt *HttpPacket) GetHost() string {
	return pkt.Host
}

// SetPort 设置请求的端口
func (pkt *HttpPacket) SetPort(value int) {
	pkt.Port = value
	//同时修正其他的相关参数
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
	pkt.SetURL(pkt.URL)
}

// GetPort 获取请求的端口
func (pkt *HttpPacket) GetPort() int {
	return pkt.Port
}

// GetHostAndPort 获取请求的Host和端口
func (pkt *HttpPacket) GetHostAndPort() string {
	var hostAndPort string = ""
	if pkt.GetPort() == 80 || pkt.GetPort() == 443 {
		hostAndPort = pkt.GetHost()
	} else {
		hostAndPort = fmt.Sprintf("%s:%d", pkt.GetHost(), pkt.GetPort())
	}
	return hostAndPort
}

// GetWebsite 获取站点
func (pkt *HttpPacket) GetWebsite() string {
	var website = url.URL{
		Host: pkt.GetHostAndPort(),
		Path: "/",
	}
	if pkt.IsSSL {
		website.Scheme = "https"
	} else {
		website.Scheme = "http"
	}
	return website.String()
}

// SetPath 设置请求的Path
func (pkt *HttpPacket) SetPath(value string) {
	pkt.Path = value
	//同时修正其他的相关参数
	if len(pkt.RawQuery) > 0 {
		if len(pkt.Path) == 0 {
			pkt.Path = "/" + pkt.Path
		}
		pkt.PathAndQuery = pkt.Path + "?" + pkt.RawQuery
	} else {
		pkt.PathAndQuery = pkt.Path
	}
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// GetPath 获取请求的Path
func (pkt *HttpPacket) GetPath() string {
	return pkt.Path
}

// AddQuery 添加请求的Query
func (pkt *HttpPacket) AddQuery(key, value string) {
	MapArrayAdd(pkt.Query, key, value)
	pkt.RawQuery = MapArrayToString(pkt.Query, "&", "=")
	//同时修正其他的相关参数
	if len(pkt.RawQuery) > 0 {
		pkt.PathAndQuery = pkt.Path + "?" + pkt.RawQuery
	} else {
		pkt.PathAndQuery = pkt.Path
	}
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// SetQuery 设置请求的Query
func (pkt *HttpPacket) SetQuery(key, value string) {
	MapArraySet(pkt.Query, key, value)
	pkt.RawQuery = MapArrayToString(pkt.Query, "&", "=")
	//同时修正其他的相关参数
	if len(pkt.RawQuery) > 0 {
		pkt.PathAndQuery = pkt.Path + "?" + pkt.RawQuery
	} else {
		pkt.PathAndQuery = pkt.Path
	}
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// DelQuery 删除请求的Query
func (pkt *HttpPacket) DelQuery(key string) {
	MapArrayDel(pkt.Query, key)
	pkt.RawQuery = MapArrayToString(pkt.Query, "&", "=")
	//同时修正其他的相关参数
	if len(pkt.RawQuery) > 0 {
		pkt.PathAndQuery = pkt.Path + "?" + pkt.RawQuery
	} else {
		pkt.PathAndQuery = pkt.Path
	}
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// GetQuery 获取请求的Query
func (pkt *HttpPacket) GetQuery(key string) string {
	return MapArrayGet(pkt.Query, key)
}

// SetRawQuery 设置请求的Query
func (pkt *HttpPacket) SetRawQuery(value string) {
	pkt.RawQuery = value
	pkt.Query = make(map[string][]string)
	tmp := strings.Split(value, "&")
	for _, v := range tmp {
		tmp2 := strings.SplitN(v, "=", 2)
		if len(tmp2) > 1 {
			MapArrayAdd(pkt.Query, tmp[0], tmp[1])
		} else {
			MapArrayAdd(pkt.Query, tmp[0], "")
		}
	}
	//同时修正其他的相关参数
	if len(pkt.RawQuery) > 0 {
		pkt.PathAndQuery = pkt.Path + "?" + pkt.RawQuery
	} else {
		pkt.PathAndQuery = pkt.Path
	}
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// GetRawQuery 获取请求的Query
func (pkt *HttpPacket) GetRawQuery() string {
	return pkt.RawQuery
}

// SetPathAndQuery 设置请求的Path和Query
func (pkt *HttpPacket) SetPathAndQuery(value string) {
	pkt.PathAndQuery = value
	tmp := strings.SplitN(pkt.PathAndQuery, "?", 2)
	if len(tmp) > 1 {
		pkt.Path = tmp[0]
		pkt.SetRawQuery(tmp[1])
	} else {
		pkt.Path = tmp[0]
		pkt.SetRawQuery("")
	}
}

// GetPathAndQuery 获取请求的Path和Query
func (pkt *HttpPacket) GetPathAndQuery() string {
	return pkt.PathAndQuery
}

// SetUserAgent 设置请求的UserAgent
func (pkt *HttpPacket) SetUserAgent(value string) {
	pkt.UserAgent = value

}

// GetUserAgent 获取请求的UserAgent
func (pkt *HttpPacket) GetUserAgent() string {
	return pkt.UserAgent
}

// SetPostData 设置请求的PostData
func (pkt *HttpPacket) SetPostData(value string) {
	pkt.PostData = value
	pkt.PostDataBytes = []byte(value)
}

// GetPostData 获取请求的PostData
func (pkt *HttpPacket) GetPostData() string {
	return pkt.PostData
}

// SetFragment 设置请求的Fragment
func (pkt *HttpPacket) SetFragment(value string) {
	pkt.Fragment = value
	//同时修正其他的相关参数
	if pkt.Port == 80 || pkt.Port == 443 {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s%s#%s", pkt.Scheme, pkt.Host, pkt.PathAndQuery, pkt.Fragment)
		}
	} else {
		if len(pkt.Fragment) < 0 {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery)
		} else {
			pkt.URL = fmt.Sprintf("%s://%s:%d%s#%s", pkt.Scheme, pkt.Host, pkt.Port, pkt.PathAndQuery, pkt.Fragment)
		}
	}
}

// GetFragment 获取请求的Fragment
func (pkt *HttpPacket) GetFragment() string {
	return pkt.Fragment
}

// AddCookie 添加请求的Cookie
func (pkt *HttpPacket) AddCookie(key, value string) {
	MapArrayAdd(pkt.Cookies, key, value)
	pkt.RawCookies = MapArrayToString(pkt.Cookies, ";", "=")
}

// SetCookie 设置请求的Cookie
func (pkt *HttpPacket) SetCookie(key, value string) {
	MapArraySet(pkt.Cookies, key, value)
	pkt.RawCookies = MapArrayToString(pkt.Cookies, ";", "=")
}

// DelCookie 删除请求的Cookie
func (pkt *HttpPacket) DelCookie(key string) {
	MapArrayDel(pkt.Cookies, key)
	pkt.RawCookies = MapArrayToString(pkt.Cookies, ";", "=")
}

// GetCookie 获取请求的Cookie
func (pkt *HttpPacket) GetCookie(key string) string {
	return MapArrayGet(pkt.Cookies, key)
}

// SetRawCookie 设置请求的RawCookie
func (pkt *HttpPacket) SetRawCookie(value string) {
	pkt.Cookies = make(map[string][]string)
	tmp := strings.Split(value, ";")
	for _, v := range tmp {
		tmp2 := strings.SplitN(v, "=", 2)
		if len(tmp2) > 1 {
			//fmt.Println(tmp2[0], tmp2[1])
			MapArrayAdd(pkt.Cookies, strings.TrimSpace(tmp2[0]), tmp2[1])
		} else {
			MapArrayAdd(pkt.Cookies, strings.TrimSpace(tmp2[0]), "")
		}
	}
	pkt.RawCookies = MapArrayToString(pkt.Cookies, ";", "=")
}

// GetRawCookie 获取请求的RawCookie
func (pkt *HttpPacket) GetRawCookie() string {
	return pkt.RawCookies
}

// AddHeader 添加请求的Header
func (pkt *HttpPacket) AddHeader(key, value string) {
	MapArrayAdd(pkt.Headers, key, value)
	pkt.RawHeaders = MapArrayToString(pkt.Headers, "\r\n", ":")
}

// SetHeader 设置请求的Header
func (pkt *HttpPacket) SetHeader(key, value string) {
	MapArraySet(pkt.Headers, key, value)
	pkt.RawHeaders = MapArrayToString(pkt.Headers, "\r\n", ":")
}

// DelHeader 删除请求的Header
func (pkt *HttpPacket) DelHeader(key string) {
	MapArrayDel(pkt.Headers, key)
	pkt.RawHeaders = MapArrayToString(pkt.Headers, "\r\n", ":")
}

// GetHeader 获取请求的Header
func (pkt *HttpPacket) GetHeader(key string) string {
	return MapArrayGet(pkt.Headers, key)
}

// GetRawHeader 获取请求的RawHeader
func (pkt *HttpPacket) GetRawHeader() string {
	return pkt.RawHeaders
}

// GetFullHeader 获取请求的完整头部
func (pkt *HttpPacket) GetFullHeader() string {
	var fullHeader = ""
	if len(pkt.GetFragment()) > 0 {
		fullHeader = fmt.Sprintf("%s %s#%s %s", pkt.GetMethod(), pkt.GetPathAndQuery(), pkt.GetFragment(), pkt.GetHttpVer())
	} else {
		fullHeader = fmt.Sprintf("%s %s %s", pkt.GetMethod(), pkt.GetPathAndQuery(), pkt.GetHttpVer())
	}
	if len(pkt.GetHost()) > 0 {
		fullHeader = fmt.Sprintf("%s\r\nHost: %s", fullHeader, pkt.GetHostAndPort())
	}
	if len(pkt.GetUserAgent()) > 0 {
		fullHeader = fmt.Sprintf("%s\r\nUser-Agent: %s", fullHeader, pkt.GetUserAgent())
	}
	if len(pkt.GetRawHeader()) > 0 {
		fullHeader = fmt.Sprintf("%s\r\n%s", fullHeader, pkt.GetRawHeader())
	}
	if len(pkt.GetRawCookie()) > 0 {
		fullHeader = fmt.Sprintf("%s\r\nCookie: %s", fullHeader, pkt.GetRawCookie())
	}
	if len(pkt.GetRawDupHeader()) > 0 {
		fullHeader = fmt.Sprintf("%s\r\n%s", fullHeader, pkt.GetRawDupHeader())
	}
	if pkt.AuthUser != "" || pkt.AuthPwd != "" {
		fullHeader = fmt.Sprintf("%s\r\nAuthorization: %s", fullHeader,
			base64.StdEncoding.EncodeToString([]byte(pkt.AuthUser+":"+pkt.AuthPwd)))
	}
	return fullHeader + pkt.BodySplit
}

// AddDupHeader 添加请求的DupHeader
func (pkt *HttpPacket) AddDupHeader(key, value string) {
	MapArrayAdd(pkt.DupHeaders, key, value)
	pkt.RawDupHeaders = MapArrayToString(pkt.DupHeaders, "\r\n", ":")
}

// SetDupHeader 设置请求的DupHeader
func (pkt *HttpPacket) SetDupHeader(key, value string) {
	MapArraySet(pkt.DupHeaders, key, value)
	pkt.RawDupHeaders = MapArrayToString(pkt.DupHeaders, "\r\n", ":")
}

// DelDupHeader 删除请求的DupHeader
func (pkt *HttpPacket) DelDupHeader(key string) {
	MapArrayDel(pkt.DupHeaders, key)
	pkt.RawDupHeaders = MapArrayToString(pkt.DupHeaders, "\r\n", ":")
}

// GetDupHeader 获取请求的DupHeader
func (pkt *HttpPacket) GetDupHeader(key string) string {
	return MapArrayGet(pkt.DupHeaders, key)
}

// GetRawDupHeader 获取请求的RawDupHeader
func (pkt *HttpPacket) GetRawDupHeader() string {
	return pkt.RawDupHeaders
}

// GetOriginPacket 获取请求的原始包
func (pkt *HttpPacket) GetOriginPacket() string {
	return pkt.OriginalPacket
}

// Reset 恢复当前的包为原始的数据包
func (pkt *HttpPacket) Reset() {
	ok, newPkt := ParserPacket(pkt.MarkedPacket, pkt.IsSSL)
	if ok {
		// todo 如果要将pkt赋值为newPkt应使用 *pkt=*newPkt
		pkt = newPkt
	}
}

// GetPacket 获取请求的当前包
func (pkt *HttpPacket) GetPacket() string {
	var newPacket = ""
	if len(pkt.GetFragment()) > 0 {
		newPacket = fmt.Sprintf("%s %s#%s %s", pkt.GetMethod(), pkt.GetPathAndQuery(), pkt.GetFragment(), pkt.GetHttpVer())
	} else {
		newPacket = fmt.Sprintf("%s %s %s", pkt.GetMethod(), pkt.GetPathAndQuery(), pkt.GetHttpVer())
	}
	if len(pkt.GetHost()) > 0 {
		newPacket = fmt.Sprintf("%s\r\nHost: %s", newPacket, pkt.GetHostAndPort())
	}
	if len(pkt.GetUserAgent()) > 0 {
		newPacket = fmt.Sprintf("%s\r\nUser-Agent: %s", newPacket, pkt.GetUserAgent())
	}
	if len(pkt.GetRawHeader()) > 0 {
		newPacket = fmt.Sprintf("%s\r\n%s", newPacket, pkt.GetRawHeader())
	}
	if len(pkt.GetRawCookie()) > 0 {
		newPacket = fmt.Sprintf("%s\r\nCookie: %s", newPacket, pkt.GetRawCookie())
	}
	if len(pkt.GetRawDupHeader()) > 0 {
		newPacket = fmt.Sprintf("%s\r\n%s", newPacket, pkt.GetRawDupHeader())
	}
	if pkt.AuthUser != "" || pkt.AuthPwd != "" {
		newPacket = fmt.Sprintf("%s\r\nAuthorization: %s", newPacket,
			base64.StdEncoding.EncodeToString([]byte(pkt.AuthUser+":"+pkt.AuthPwd)))
	}
	if len(pkt.GetPostData()) > 0 {
		newPacket = fmt.Sprintf("%s%s%s", newPacket, pkt.BodySplit, pkt.GetPostData())
	} else {
		newPacket = fmt.Sprintf("%s\r\n\r\n", newPacket)
	}
	pkt.CurrentPayLoadPacket = newPacket

	return pkt.CurrentPayLoadPacket
}

// CreatePacket 构造一个数据包
func (pkt *HttpPacket) CreatePacket(packet string, ssl bool) *HttpPacket {
	ok, newPkt := ParserPacket(packet, ssl)
	if ok {
		pkt.Valid = true
	} else {
		pkt.Valid = false
	}
	return newPkt
}

// ParserPacket 分析截获的数据包
func ParserPacket(packet string, ssl bool) (bool, *HttpPacket) {
	//为空不继续
	var pkt HttpPacket
	if len(packet) == 0 {
		pkt.Valid = false
		return false, &pkt
	}
	// 原始包
	pkt.Query = make(map[string][]string)
	pkt.Cookies = make(map[string][]string)
	pkt.Headers = make(map[string][]string)
	pkt.DupHeaders = make(map[string][]string)

	pkt.MarkedPacket = packet

	if strings.Contains(packet, "§") {
		pkt.OriginalPacket = DeleteTestTag(packet)
		pkt.CurrentPayLoadPacket = packet
	} else {
		pkt.OriginalPacket = packet
	}

	//是否HTTPS
	pkt.IsSSL = ssl

	//找Body和Header的分隔符
	regex, err := regexp2.Compile(`\r*\n\r*\n`, regexp2.RE2)
	if err != nil {
		return false, &pkt
	}
	match, err := regex.FindStringMatch(packet)
	if err != nil || match == nil {
		return false, &pkt
	}

	pkt.BodySplit = match.String()

	var (
		rawHeader      string
		rawBody        string
		isGotFirstLine bool
		splitLines     []string
		pathAndQuery   string
		latestHeader   string
	)

	PacketArray := strings.SplitN(packet, pkt.BodySplit, 2)
	if len(PacketArray) < 2 {
		return false, &pkt
	} else {
		rawHeader = PacketArray[0]
		rawBody = PacketArray[1]
	}
	//将头部中的所有\r替换为空
	rawHeader = strings.ReplaceAll(rawHeader, "\r", "")
	//按行分割
	splitLines = strings.Split(rawHeader, "\n")
	for _, line := range splitLines {
		lineLower := strings.ToLower(line)
		if !isGotFirstLine { //如果没获取到第一行
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			isGotFirstLine = true
			blankSplit := strings.Split(line, " ")
			if len(blankSplit) > 2 { //GET /xxx HTTP/1.1，允许GET /xxx xxx.do xxx/xxx
				pkt.Method = blankSplit[0] //Method
				for i := 1; i <= len(blankSplit)-2; i++ {
					if len(pathAndQuery) > 0 {
						pathAndQuery = pathAndQuery + " " + blankSplit[i] //允许编辑器里出现空格
					} else {
						pathAndQuery = blankSplit[i]
					}
				}
				pkt.PathAndQuery = pathAndQuery
				pkt.HttpVer = blankSplit[len(blankSplit)-1]
			} else if len(blankSplit) == 2 { //GET /xxx,无协议
				pkt.Method = blankSplit[0] //Method
				for i := 1; i <= len(blankSplit)-1; i++ {
					if len(pathAndQuery) > 0 {
						pathAndQuery = pathAndQuery + " " + blankSplit[i] //允许编辑器里出现空格
					} else {
						pathAndQuery = blankSplit[i]
					}
				}
				pkt.PathAndQuery = pathAndQuery
				pkt.HttpVer = ""
			}
			continue
		}
		//头部未结束、新行不以空格、TAB开头,并且该行不能为空
		if len(line) > 0 && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			latestHeader = strings.Split(line, ":")[0] //最近的头部字段
			if strings.HasPrefix(lineLower, "host:") {
				host := str.StringReplaceWithIgnoreCase(line, "host:", "")
				host = strings.TrimSpace(host)
				//获取端口
				if strings.Contains(host, ":") {
					portStr := strings.Split(host, ":")[1]
					if str.IsNumeric(portStr) {
						port := str.String2IntWithDefaultValue(portStr, 0)
						if port > 0 && port <= 65535 {
							pkt.Port = port
						} else {
							pkt.Port = 80
						}
					} else {
						pkt.Port = 80
					}
				} else {
					if ssl {
						pkt.Port = 443
					} else {
						pkt.Port = 80
					}
				}
				//获取主机
				h := strings.Split(host, ":")[0]
				if strings.Contains(h, "§") {
					h = strings.Replace(h, "§", "", -1)
				}
				pkt.Host = h
				if pkt.Port != 80 && pkt.Port != 443 {
					pkt.HostAndPort = pkt.Host + ":" + strconv.Itoa(pkt.Port)
				} else {
					pkt.HostAndPort = pkt.Host
				}
				//获取URL
				if strings.HasPrefix(pkt.PathAndQuery, "/") {
					pkt.URL = pkt.HostAndPort + pkt.PathAndQuery
				} else if strings.HasPrefix(strings.ToLower(pkt.PathAndQuery), "http://") || strings.HasPrefix(strings.ToLower(pkt.PathAndQuery), "https://") {
					pkt.URL = pkt.PathAndQuery
				} else {
					pkt.URL = pkt.HostAndPort
				}
				if !strings.HasPrefix(strings.ToLower(pkt.PathAndQuery), "http://") && !strings.HasPrefix(strings.ToLower(pkt.PathAndQuery), "https://") {
					if ssl {
						pkt.URL = "https://" + pkt.URL
					} else {
						pkt.URL = "http://" + pkt.URL
					}
				}
			} else if strings.HasPrefix(lineLower, "user-agent:") { //获取用户代理
				if len(pkt.GetUserAgent()) == 0 {
					pkt.UserAgent = strings.TrimSpace(str.StringReplaceWithIgnoreCase(line, "user-agent:", ""))
				} else {
					pkt.AddDupHeader("User-Agent", strings.TrimSpace(str.StringReplaceWithIgnoreCase(line, "user-agent:", "")))
				}
			} else if strings.HasPrefix(lineLower, "cookie:") { //获取Cookie
				if len(pkt.GetRawCookie()) == 0 {
					pkt.SetRawCookie(strings.TrimSpace(str.StringReplaceWithIgnoreCase(line, "cookie:", "")))
				} else {
					pkt.AddDupHeader("Cookie", strings.TrimSpace(str.StringReplaceWithIgnoreCase(line, "cookie:", "")))
				}
			} else if strings.HasPrefix(lineLower, "accept-encoding:") { //不需要Accept-Encoding:,部分站点不支持gzip，会乱码
				continue
				//} else if strings.HasPrefix(lineLower, "content-length:") { //不需要Content-Length
				//	continue
			} else if strings.HasPrefix(lineLower, "connection:") { //不需要Connection
				continue
			} else {
				//获取其他头部
				tmp := strings.SplitN(line, ":", 2)
				if _, ok := pkt.Headers[tmp[0]]; !ok {
					if len(tmp) > 1 {
						pkt.AddHeader(tmp[0], tmp[1]) //值可能为空
					} else {
						pkt.AddHeader(tmp[0], "§NIL§") //没有值
					}
				} else {
					if len(tmp) > 1 {
						pkt.AddDupHeader(tmp[0], tmp[1]) //值可能为空
					} else {
						pkt.AddDupHeader(tmp[0], "§NIL§") //没有值
					}
				}
			}
		} else if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") { //头部未结束、新行以空格或者tab开头
			//在最后一个头部追加数据
			if len(latestHeader) > 0 {
				if strings.TrimSpace(strings.ToLower(latestHeader)) == "user-agent" {
					pkt.UserAgent = pkt.UserAgent + "\r\n" + line
				} else if strings.TrimSpace(strings.ToLower(latestHeader)) == "cookie" {
					pkt.SetRawCookie(fmt.Sprintf("%s%s%s", pkt.RawCookies, "\r\n", line))
				} else {
					pkt.SetHeader(latestHeader, fmt.Sprintf("%s%s%s", pkt.GetHeader(latestHeader), "\r\n", line))
				}
			}
		}
	}
	if len(rawBody) > 0 {
		pkt.PostData = rawBody
		pkt.PostDataBytes = []byte(pkt.PostData)
	} else {
		pkt.PostData = ""
	}

	pkt.DestIP = pkt.Host
	pkt.DestPort = pkt.Port

	pkt.URL = strings.ReplaceAll(pkt.URL, "[", "%5B")
	pkt.URL = strings.ReplaceAll(pkt.URL, "]", "%5D")

	uri, err := url.Parse(pkt.URL)
	if err != nil {
		pkt.Valid = false
		return false, &pkt
	}
	pkt.Scheme = uri.Scheme
	pkt.Path = uri.Path
	pkt.RawQuery = uri.RawQuery
	pkt.Query = uri.Query()

	pkt.Valid = true

	return true, &pkt
}

// DeleteTestTag 清除所有变体标记
func DeleteTestTag(packet string) string {
	return strings.Replace(packet, "§", "", -1)
}
