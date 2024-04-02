package iconvCharset

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/google/btree"
	"gopkg.in/iconv.v1"
)

var iconvSupportCharset *btree.BTree

var priorCharset []Charset

func init() {
	priorCharset = []Charset{"UTF-8", "GB2312", "GBK", "GB18030", "BIG5"}
	iconvSupportCharset = btree.New(3)
	br := bufio.NewReader(bytes.NewReader([]byte(charset)))
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if len(line) > 0 {
			iconvSupportCharset.ReplaceOrInsert(Charset(line))
		}
	}
}

type Charset string

func (c Charset) Less(than btree.Item) bool {
	return c < than.(Charset)
}

func AllSupport() []Charset {
	var res []Charset
	res = append(res, priorCharset...)
	iconvSupportCharset.Ascend(func(i btree.Item) bool {
		for _, item := range priorCharset {
			if i.(Charset) == item {
				return true
			}
		}
		res = append(res, i.(Charset))
		return true
	})
	return res
}

func IsSupport(_charset Charset) bool {
	return iconvSupportCharset.Has(_charset)
}

func ConvertTo(dst Charset, src Charset, input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("input is null")
	}
	if !IsSupport(dst) {
		return nil, fmt.Errorf("iconv iconvCharset not support this: %v", dst)
	}
	if !IsSupport(src) {
		return nil, fmt.Errorf("iconv iconvCharset not support this: %v", src)
	}
	cd, err := iconv.Open(string(dst), string(src))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(iconv.NewReader(cd, bytes.NewReader(input), 0))
}
