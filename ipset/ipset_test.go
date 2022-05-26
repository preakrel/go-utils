package ipset

import (
	"testing"
)

func TestAddOne(t *testing.T) {
	data := [][]string{
		{"0.0.0.0", "0.0.0.1"},
		{"1.255.254.255", "1.255.255.0"},
		{"1.255.255.255", "2.0.0.0"},
		{"255.255.255.255", "0.0.0.0"},
		{"::", "::1"},
		{"::ffff", "::1:0"},
		{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "::"},
	}
	for _, item := range data {
		in := ParseIp(item[0])
		expect := ParseIp(item[1])
		result := addOne(in)
		assertEqualsF(expect, result, "expect %s + 1 => %s, but got %s", in, expect, result)
	}
}

func TestIpNetWrapper_String(t *testing.T) {
	inputs := []string{
		"0.0.0.0/0",
		"0.0.0.0/1",
		"0.0.0.0/32",
		"128.0.0.0/1",
		"0.0.1.0/24",
		"192.0.1.0/24",
		"::/0",
	}
	for _, text := range inputs {
		iRange, err := Parse(text)
		assert(err == nil, "err is %v", err)
		ipNetWrapper := iRange.(IpNetWrapper)
		assertEquals(text, ipNetWrapper.String())
	}
}
func TestIpMerge(t *testing.T) {
	inputs := []string{
		"1.9.0.0", "1.9.7.255",
		"1.9.8.0", "1.9.8.63",
		"1.9.8.64", "1.9.8.255",
	}
	data := make([]IRange, 0, len(inputs))
	for _, text := range inputs {
		iRange, err := Parse(text)
		assert(err == nil, "err is %v", err)
		data = append(data, iRange)

	}
	//处理结果（掩码合并为IP段）
	result := ConvertBatch(SortAndMerge(data), 2)
	for _, v := range result {
		t.Log(ReturnSelf(v))
	}
}
