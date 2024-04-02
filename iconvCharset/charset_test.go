package iconvCharset

import (
	"testing"
)

func TestAll(t *testing.T) {
	t.Log(AllSupport())
}

func TestConvertTo(t *testing.T) {
	//res, err := ConvertTo("UTF8", "WINDOWS-1252", []byte("1234"))
	//t.Log(string(res), err)
	res, err := ConvertTo("GBK", "UTF8", []byte("哈哈哈"))
	t.Log(string(res), err)
	res, err = ConvertTo("UTF8", "GBK", res)
	t.Log(string(res), err)
}
