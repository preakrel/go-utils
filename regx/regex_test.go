package regx

import (
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

var packet = `POST /admin/login HTTP/1.1
Host: elm.cangdu.org:§80§
ProxyID-Connection: keep-alive
Cache-Control: no-cache
Accept-Language: en-US,en;q=0.9
Dnt: 1
Pragma: no-cache
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36
Accept: */*
Content-Type: application/json
Cookie:asdasd
Referer: https://cangdu.org/

{"user_name":"§123§","password":"§456§"}`

func TestRegexMatch(t *testing.T) {
	t.Log(RegexMatch("12345aAA6789", "aaa"))
	t.Log(RegexMatch("12345aAA6789", "(?msi)aaa"))
	t.Log(RegexMatch(`{"transactions": [{"transactionld": "1", "amount":10000, "currency": “RON", "from": "John Doe", "to": "Dane Doe", "details": "Go shopping!"), {"transactionld": "2", "amount":500, "currency": "RON", "from": "John Doe", "to": "Unknown Company", "details": "Monthly bill 10.01.2017"}]}`, `(?msi)^{[\s\S]*}$`))
}

func TestStringToRegexEscape(t *testing.T) {
	t.Log(StringToRegexEscape(`(a<b1234$%\/:xafa)`))
}

func TestStringReplaceWithIgnoreCase(t *testing.T) {
	t.Log(string(StringReplaceWithIgnoreCase("12345aAA6789", "aaa", "test")))
	t.Log(string(StringReplaceWithIgnoreCase("12345aAA6789", "^1234:", "test")))
	t.Log(string(StringReplaceWithIgnoreCase("^12345aAA6789", "^1234", "test")))
	t.Log(strings.TrimSpace(StringReplaceWithIgnoreCase("COokie: xxx=123", "cookie:", "")))

}

func TestRegexReplace(t *testing.T) {
	t.Log(string(RegexReplace("12345aAA6789", "aaa", "test", regexp2.IgnoreCase)))
	t.Log(string(RegexReplace("12345aAA6789", "(?msi)aaa", "test", regexp2.Compiled)))
	t.Log(string(RegexReplace("12345aAA6789", "aaa", "test", regexp2.Compiled)))
	t.Log(string(RegexReplace("12345aAA6789", "^12345A", "test", regexp2.Compiled)))
	t.Log(string(RegexReplace("12345aAA6789", "(?msi)^12345A", "test", regexp2.Compiled)))
	t.Log(string(RegexReplace("12345aAA6789", "^12345A", "test", regexp2.IgnoreCase)))
}

func TestExtractWithRegex(t *testing.T) {
	s, _ := ExtractWithRegex(packet, "§(.*?)§", true)
	t.Log(s)
	s, _ = ExtractWithRegex(packet, `(?msi)Cookie:(.*?)\r?\n?$`, true)
	t.Log(s)
	s, _ = ExtractWithRegex(packet, `(?msi)Host:\s*(.*?)\r?\n`, false)
	t.Log(s)
	boundary, _ := ExtractWithRegex(`multipart/form-data; boundary=--------720221776`, "(?<=boundary\\=).*(?=$)", true)
	t.Log(boundary) //--------720221776
	r := multipart.NewReader(strings.NewReader(`----------720221776
Content-Disposition: form-data; name="addr"

xxx
----------720221776
Content-Disposition: form-data; name="count"

5
----------720221776
Content-Disposition: form-data; name="start"

Start
----------720221776
Content-Disposition: form-data; name="result"


----------720221776--`), boundary)
	for {
		part, err := r.NextRawPart()
		if err != nil {
			t.Log(err)
			return
		}
		t.Log(part.FormName())
		b, e := io.ReadAll(part)
		t.Log(string(b), e)
	}
}

func TestExtractWithRegex2(t *testing.T) {
	t.Log(string(ExtractWithRegex2(packet, "§(.*?)§", regexp2.IgnoreCase)))
	t.Log(string(ExtractWithRegex2(packet, `(?msi)Cookie:(.*?)\r?\n?$`, regexp2.IgnoreCase)))
	t.Log(string(ExtractWithRegex2(packet, `(?msi)Host:\s*(.*?)\r?\n`, regexp2.IgnorePatternWhitespace)))

}

func TestExtractSub_With_Regex(t *testing.T) {
	re := `(?mi)host:\s*(.*?)\r?\n?$`
	p := "Host: elm.cangdu.org:§80§"
	h, _ := ExtractWithRegex(p, re, false)
	t.Log(string(h))
	p = "Host:elm.cangdu.org:§80§\r\n"
	h, _ = ExtractWithRegex(p, re, false)
	t.Log(string(h))
	p = "Host:    elm.cangdu.org:§80§\r\n"
	h, _ = ExtractWithRegex(p, re, false)
	t.Log(string(h))
	p = "Host:\t\r\nelm.cangdu.org:§80§\r\n"
	h, _ = ExtractWithRegex(p, re, false)
	t.Log(string(h))
}

func TestExtractSilceWithRegex(t *testing.T) {
	t.Log(ExtractSliceWithRegex(`@-moz-keyframes yy-ipt{0%{border-color:#4791ff transparent #4791ff #4791ff}
100%{border-color:#e10602 transparent #e10602 #e10602}}
@-webkit-keyframes yy-ipt{0%{border-color:#4791ff transparent #4791ff #4791ff}
100%{border-color:#e10602 transparent #e10602 #e10602}}
@-o-keyframes yy-ipt{0%{border-color:#4791ff transparent #4791ff #4791ff}
100%{border-color:#e10602 transparent #e10602 #e10602}}
@keyframes yy-ipt-blue{0%{border-color:#e10602 transparent #e10602 #e10602}
100%{border-color:#4791ff transparent #4791ff #4791ff}}
@-moz-keyframes yy-ipt-blue{0%{border-color:#e10602 transparent #e10602 #e10602}
100%{border-color:#4791ff transparent #4791ff #4791ff}}`, "(?msi)(?<=border-color:).*?(?= )", regexp2.Compiled))

	t.Log(ExtractSliceWithRegex(`----------720221776
Content-Disposition: form-data; name="§addr§"

§xxx§
----------720221776
Content-Disposition: form-data; name="§count§"

§5§
----------720221776
Content-Disposition: form-data; name="§start§"

§Start§
----------720221776
Content-Disposition: form-data; name="§result§"

§§
----------720221776--`, `(?<=form-data; name\=("|'))[\S\s]+?§.*?§"\r?\n`, regexp2.Compiled))
}
