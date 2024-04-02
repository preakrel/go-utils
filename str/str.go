package str

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"util/iconvCharset"

	"github.com/dlclark/regexp2"
	"github.com/gogs/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/iconv.v1"
)

// RegexMatch 正则表达式匹配
func RegexMatch(s, regx string) bool {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	//如果是包含分组的正则表达式
	r, err := regexp2.Compile(`\(\?P\<\w+\>`, regexp2.Compiled)
	if err != nil {
		return false
	}
	b, err := r.MatchString(regx)
	if err != nil {
		return false
	}
	if b {
		re, err := regexp.Compile(regx)
		if err != nil {
			return false
		}
		return re.MatchString(s)
	}
	r, err = regexp2.Compile(regx, regexp2.Compiled)
	if err != nil {
		return false
	}
	b, err = r.MatchString(s)
	if err != nil {
		return false
	}
	return b
}

// RegexMatchWithIgnoreCase 正则表达式匹配
func RegexMatchWithIgnoreCase(s, regx string) (bool, error) {
	//如果是包含分组的正则表达式
	if RegexMatch(regx, `\(\?P\<\w+\>`) {
		re, err := regexp.Compile(regx)
		if err != nil {
			return false, err
		}
		return re.MatchString(s), nil
	}
	if !RegexMatch(regx, `^\(\?[msi]{1,3}\)`) {
		regx = "(?msi)" + regx
	}
	r, err := regexp2.Compile(regx, regexp2.Compiled)
	if err != nil {
		return false, err
	}
	b, err := r.MatchString(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// GetCharsetName 获取编码
func GetCharsetName(s string, isCode bool) string {
	var resultCharset = ""
	var options = regexp2.RE2 | regexp2.IgnoreCase | regexp2.Multiline | regexp2.Singleline
	if isCode {
		resultCharset = ExtractWithRegex2(s, `(?msi)<Meta[^>]*?resultCharset=["']?([-a-zA-Z_0-9]+)["']?`, regexp2.RegexOptions(options))
	} else {
		resultCharset = ExtractWithRegex2(s, `(?msi)resultCharset=([-a-zA-Z_0-9]+)`, regexp2.RegexOptions(options))
	}
	if len(resultCharset) > 0 {
		resultCharset = strings.Replace(strings.ToLower(resultCharset), "resultCharset=", "", -1)
		resultCharset = strings.Replace(resultCharset, "'", "", -1)
		resultCharset = strings.Replace(resultCharset, "\"", "", -1)
		resultCharset = strings.Replace(resultCharset, "’", "", -1)
		resultCharset = strings.Replace(resultCharset, "‘", "", -1)
		resultCharset = strings.Replace(resultCharset, "“", "", -1)
		resultCharset = strings.Replace(resultCharset, "”'", "", -1)
		resultCharset = strings.Replace(resultCharset, " ", "", -1)
	}
	if resultCharset == "utf8" {
		resultCharset = "utf-8"
	}
	return resultCharset
}

// ExtractWithRegex2 正则表达式提取2
func ExtractWithRegex2(s, r string, options regexp2.RegexOptions) string {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	regex, err := regexp2.Compile(r, options)
	if err != nil {
		return ""
	}
	match, err := regex.FindStringMatch(s)
	if err != nil || match == nil {
		return ""
	}
	return match.String()
}

func ExtractWithRegexMatch(s, r string, options regexp2.RegexOptions) *regexp2.Match {
	regex, err := regexp2.Compile(r, options)
	if err != nil {
		return nil
	}
	match, err := regex.FindStringMatch(s)
	if err != nil || match == nil {
		return nil
	}
	return match
}

// Convert 编码转换将src编码转换为dst编码
// https://www.gnu.org/software/libiconv/
func Convert(dst, src string, reader io.Reader) ([]byte, error) {
	const (
		// MaxResponseBytes 读取返回包最大字节数
		MaxResponseBytes = 20971520 //20MB
	)
	if dst == "" {
		dst = "utf-8"
	}
	if src == "" {
		src = "utf-8"
	}
	if dst == src {
		return io.ReadAll(reader)
	}
	cd, err := iconv.Open(dst, src) // convert src to dst
	if err != nil {
		return nil, err
	}
	defer cd.Close()
	r := iconv.NewReader(cd, reader, 0)
	var buffer = make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			io.Copy(io.Discard, r) //释放掉未获取的数据
			return nil, err
		}
		if n == 0 {
			break
		}
		buffer = append(buffer, buf[:n]...)
		if len(buffer) > MaxResponseBytes {
			io.Copy(io.Discard, r) //释放掉未获取的数据
			break
		}
	}
	return buffer, nil
}

type charset string

const (
	UTF8    = charset("UTF-8")
	GB18030 = charset("GB18030")
)

// ConvertByte 将b charset转换为UTF-8编码
func ConvertByte(b []byte, charset charset) []byte {
	var byt = make([]byte, 0)
	switch charset {
	case GB18030:
		byt, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(b)
	case UTF8:
		fallthrough
	default:
		byt = b
	}
	return byt
}

func EncodeSha256(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum([]byte("")))
}

type CmdHelpInfo struct {
	ParamName    string `json:"param_name"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`

	leftLineSpaceNumber int //行前空格数 用来判断换行是否是当前的描述
}

// ParseHelpText 解析cmd的help帮助信息
func ParseHelpText(helpText string) (options []*CmdHelpInfo, err error) {
	optionRegex := regexp2.MustCompile(`(?m)^\s*(-[^,\s]+(?:,\s*--?[^\s,]+)?)(?:\s*=\s*(\S+))?\s*(.*)$`, 0)
	var regexOptions = regexp2.RE2 | regexp2.IgnoreCase | regexp2.Multiline | regexp2.Singleline
	lines := strings.Split(strings.ReplaceAll(helpText, "\r\n", "\n"), "\n")
	for _, line := range lines {
		if matches, err := optionRegex.FindStringMatch(line); err == nil {
			desc := line //默认当前行为描述
			leftLineSpaceNumber := len(line) - len(strings.TrimLeft(line, " "))
			defaultVal := ""
			if matches != nil {
				// 输出匹配结果
				for matches != nil {
					mg := matches.Groups()
					mLen := len(mg)
					var pName string
					if mLen > 1 {
						pName = strings.TrimRight(mg[1].String(), ":")
					}
					if mLen > 3 {
						desc = mg[3].String()
						if ewg := ExtractWithRegexMatch(desc, `(?msi)\(.*?default(?:[:\s]+(.*?))?\)`, regexp2.RegexOptions(regexOptions)); ewg != nil {
							defaultVal = strings.TrimSpace(ewg.GroupByNumber(1).String())
						}
					}

					if pName != "" {
						for _, vName := range strings.Split(pName, ",") {
							if vNames := strings.Split(vName, "="); len(vNames) > 1 {
								vName, desc = strings.TrimSpace(vNames[0]), fmt.Sprintf("%s%s", strings.Join(vNames[1:], ""), desc)
							}
							option := &CmdHelpInfo{ParamName: strings.TrimSpace(vName), DefaultValue: defaultVal, Description: desc, leftLineSpaceNumber: leftLineSpaceNumber}
							options = append(options, option)
						}
					}
					matches, _ = optionRegex.FindNextMatch(matches)
				}
			} else {
				//换行且没有以-开头 则是描述
				if len(options) > 0 {
					preOption := options[len(options)-1]
					//前面空格数大于前面的，则说明是描述换行
					if preOption.leftLineSpaceNumber < leftLineSpaceNumber {
						if preOption.DefaultValue == "" {
							if ewg := ExtractWithRegexMatch(line, `(?msi)\(.*?default(?:[:\s]+(.*?))?\)`, regexp2.RegexOptions(regexOptions)); ewg != nil {
								defaultVal = strings.TrimSpace(ewg.GroupByNumber(1).String())
							}
						}
						preOption.Description = fmt.Sprintf("%s %s", preOption.Description, strings.TrimSpace(line))
						if defaultVal != "" {
							preOption.DefaultValue = defaultVal
						}
						options[len(options)-1] = preOption
					}
				}
			}
		}
	}
	return
}

// CharsetConvertToUTF8 将源编码转换为UTF-8
func CharsetConvertToUTF8(input []byte, dstCharset string) []byte {
	if input == nil {
		return input
	}
	if dstCharset == "" || dstCharset == "noop" { //自动检测
		res, err := chardet.NewTextDetector().DetectBest(input)
		if err != nil {
			//log.Printf("自动检测编码失败：%s", err.Error())
			return input
		}
		dstCharset = res.Charset
	}
	dstCharset = strings.ToUpper(dstCharset)
	if dstCharset != "" && iconvCharset.IsSupport(iconvCharset.Charset(dstCharset)) {
		convBytes, err := iconvCharset.ConvertTo("UTF-8", iconvCharset.Charset(dstCharset), input)
		if err != nil {
			//log.Printf("编码[%s]转换成UTF-8失败：%s", strings.ToUpper(dstCharset), err.Error())
			return input
		}
		return convBytes
	}
	return input
}

// CharsetConvert 指定srcCharset 编码转换dstCharset编码
func CharsetConvert(input []byte, dstCharset, srcCharset string) []byte {
	if input == nil {
		return input
	}
	if srcCharset == "" || srcCharset == "noop" { //自动检测
		res, err := chardet.NewTextDetector().DetectBest(input)
		if err != nil {
			srcCharset = "UTF-8"
		} else {
			srcCharset = res.Charset
		}
	}
	srcCharset = strings.ToUpper(srcCharset)
	if srcCharset != "" && iconvCharset.IsSupport(iconvCharset.Charset(dstCharset)) && iconvCharset.IsSupport(iconvCharset.Charset(srcCharset)) {
		convBytes, err := iconvCharset.ConvertTo(iconvCharset.Charset(dstCharset), iconvCharset.Charset(srcCharset), input)
		if err != nil {
			return input
		}
		return convBytes
	}
	return input
}

func FormatJSON(jsonBts []byte) ([]byte, error) {
	var data interface{}
	err := json.Unmarshal(jsonBts, &data)
	if err != nil {
		return nil, err
	}
	formattedJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
}

func UnFormatJSON(jsonBts []byte) ([]byte, error) {
	var data interface{}
	err := json.Unmarshal(jsonBts, &data)
	if err != nil {
		return nil, err
	}
	formattedJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
}
