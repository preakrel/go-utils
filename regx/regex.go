package regx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

// RegexMatch 正则表达式匹配
func RegexMatch(s, register string) bool {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	//如果是包含分组的正则表达式
	r, err := regexp2.Compile(`\(\?P\<\w+\>`, regexp2.Compiled)
	if err != nil {
		return false
	}
	b, err := r.MatchString(register)
	if err != nil {
		return false
	}
	if b {
		re, err := regexp.Compile(register)
		if err != nil {
			return false
		}
		return re.MatchString(s)
	}
	r, err = regexp2.Compile(register, regexp2.Compiled)
	if err != nil {
		return false
	}
	b, err = r.MatchString(s)
	if err != nil {
		return false
	}
	return b
}

// RegexMatch 正则表达式匹配
func RegexMatchWithIgnoreCase(s, register string) (bool, error) {
	//如果是包含分组的正则表达式
	if RegexMatch(register, `\(\?P\<\w+\>`) {
		re, err := regexp.Compile(register)
		if err != nil {
			return false, err
		}
		return re.MatchString(s), nil
	}
	if !RegexMatch(register, `^\(\?[msi]{1,3}\)`) {
		register = "(?msi)" + register
	}
	r, err := regexp2.Compile(register, regexp2.Compiled)
	if err != nil {
		return false, err
	}
	b, err := r.MatchString(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

// ContainsEx 正则表达式实现的不区分大小写的Contains
func ContainsEx(s, subStr string) bool {
	register := StringToRegexEscape(subStr)
	r, err := regexp2.Compile(register, regexp2.Compiled)
	if err != nil {
		return false
	}
	b, err := r.MatchString(s)
	if err != nil {
		return false
	}
	return b
}

// StringRegexMatch2 正则表达式匹配2
func StringRegexMatch2(s, regStr string, ignoreCase bool) bool {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	var options = regexp2.RE2
	if ignoreCase {
		options = regexp2.RE2 | regexp2.IgnoreCase | regexp2.Multiline | regexp2.Singleline
	}
	reg, err := regexp2.Compile(regStr, regexp2.RegexOptions(options))
	if err != nil {
		return false
	}
	match, _ := reg.MatchString(s)
	return match
}

// ExtractWithRegex 正则表达式提取
func ExtractWithRegex(s, r string, ignoreCase bool) (string, error) {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	//如果是包含分组的正则表达式
	if RegexMatch(r, `\(\?P\<\w+\>`) {
		re, e := regexp.Compile(r)
		if e != nil {
			return "", e
		}
		match := re.FindStringSubmatch(s)
		if match == nil {
			return "", fmt.Errorf("正则表达式未能匹配到结果")
		}
		groupNames := re.SubexpNames()
		if len(groupNames) > 1 {
			return match[1], nil
		} else {
			return match[0], nil
		}
	}
	if !RegexMatch(r, `^\(\?[msi]{1,3}\)`) {
		if ignoreCase {
			r = "(?msi)" + r
		}
	}
	regex, err := regexp2.Compile(r, regexp2.RE2)
	if err != nil {
		return "", err
	}
	match, err := regex.FindStringMatch(s)
	if err != nil {
		return "", err
	}
	if match == nil {
		return "", fmt.Errorf("正则表达式未能匹配到结果")
	}
	return match.String(), nil
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

// ExtractSliceWithRegex 正则表达式提取2
func ExtractSliceWithRegex(s, r string, options regexp2.RegexOptions) []string {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	array := make([]string, 0)
	regex, err := regexp2.Compile(r, options)
	if err != nil {
		return array
	}
	match, err := regex.FindStringMatch(s)
	if err != nil || match == nil {
		return array
	}
	array = append(array, match.String())
	for {
		match, err = regex.FindNextMatch(match)
		if err != nil || match == nil {
			break
		}
		array = append(array, match.String())
	}
	return array
}

// RegexReplace 正则表达式替换
func RegexReplace(s, register, replacement string, options regexp2.RegexOptions) string {
	r, err := regexp2.Compile(register, options)
	if err != nil {
		return s
	}
	result, err := r.Replace(s, replacement, 0, -1)
	if err != nil {
		return s
	}
	return result
}

// RegexReplaceBackErr 正则表达式替换,返回异常
func RegexReplaceBackErr(s, expr, replacement string, options regexp2.RegexOptions) (string, error) {
	r, err := regexp2.Compile(expr, options)
	if err != nil {
		return s, err
	}
	return r.Replace(s, replacement, 0, -1)
}

// StringReplaceWithIgnoreCase 普通含特殊字符的字符串忽略大小写替换，正则表达式实现
func StringReplaceWithIgnoreCase(s, oldStr, newStr string) string {
	oldStr = strings.ReplaceAll(oldStr, "=", `\=`)
	oldStr = strings.ReplaceAll(oldStr, `!`, `\!`)
	oldStr = strings.ReplaceAll(oldStr, `.`, `\.`)
	oldStr = strings.ReplaceAll(oldStr, `|`, `\|`)
	oldStr = strings.ReplaceAll(oldStr, `\`, `\\`)
	oldStr = strings.ReplaceAll(oldStr, `*`, `\*`)
	oldStr = strings.ReplaceAll(oldStr, `?`, `\?`)
	oldStr = strings.ReplaceAll(oldStr, `<`, `\<`)
	oldStr = strings.ReplaceAll(oldStr, `>`, `\>`)
	oldStr = strings.ReplaceAll(oldStr, `+`, `\+`)
	oldStr = strings.ReplaceAll(oldStr, `{`, `\{`)
	oldStr = strings.ReplaceAll(oldStr, `}`, `\}`)
	oldStr = strings.ReplaceAll(oldStr, `(`, `\(`)
	oldStr = strings.ReplaceAll(oldStr, `)`, `\)`)
	oldStr = strings.ReplaceAll(oldStr, `:`, `\:`)
	oldStr = strings.ReplaceAll(oldStr, `^`, `\^`)
	oldStr = strings.ReplaceAll(oldStr, `$`, `\$`)
	oldStr = strings.ReplaceAll(oldStr, `[`, `\[`)
	oldStr = strings.ReplaceAll(oldStr, `]`, `\]`)
	var result = ""
	var err error
	reg, err := regexp2.Compile(oldStr, regexp2.IgnoreCase)
	if err != nil {
		return s
	}
	result, err = reg.Replace(s, newStr, 0, -1)
	if err != nil {
		return s
	}
	return result
}

// StringToRegexEscape 将字符串转换成正则表达式
func StringToRegexEscape(register string) string {
	if len(register) == 0 {
		return ""
	}
	register = strings.ReplaceAll(register, "\r", `\r`)
	register = strings.ReplaceAll(register, "\n", `\n`)
	register = strings.ReplaceAll(register, "=", `\=`)
	register = strings.ReplaceAll(register, `.`, `\.`)
	register = strings.ReplaceAll(register, `!`, `\!`)
	register = strings.ReplaceAll(register, `,`, `\,`)
	register = strings.ReplaceAll(register, `-`, `\-`)
	register = strings.ReplaceAll(register, `|`, `\|`)
	register = strings.ReplaceAll(register, `\`, `\\`)
	register = strings.ReplaceAll(register, `*`, `\*`)
	register = strings.ReplaceAll(register, `?`, `\?`)
	register = strings.ReplaceAll(register, `<`, `\<`)
	register = strings.ReplaceAll(register, `>`, `\>`)
	register = strings.ReplaceAll(register, `+`, `\+`)
	register = strings.ReplaceAll(register, `{`, `\{`)
	register = strings.ReplaceAll(register, `}`, `\}`)
	register = strings.ReplaceAll(register, `(`, `\(`)
	register = strings.ReplaceAll(register, `)`, `\)`)
	register = strings.ReplaceAll(register, `:`, `\:`)
	register = strings.ReplaceAll(register, `^`, `\^`)
	register = strings.ReplaceAll(register, `&`, `\&`)
	register = strings.ReplaceAll(register, `$`, `\$`)
	register = strings.ReplaceAll(register, `%`, `\%`)
	register = strings.ReplaceAll(register, `#`, `\#`)
	register = strings.ReplaceAll(register, `~`, `\~`)
	register = strings.ReplaceAll(register, `[`, `\[`)
	register = strings.ReplaceAll(register, `]`, `\]`)
	return register
}

//// ExtractSubRegex 正则表达式提取
//func ExtractSubRegex(s, expr string) string {
//	regex, err := regexp.Compile(expr)
//	if err != nil {
//		return ""
//	}
//	result := regex.FindStringSubmatch(s)
//	l := len(result)
//	if l > 0 {
//		return result[l-1]
//	} else {
//		return ""
//	}
//}
