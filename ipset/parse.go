package ipset

import (
	"net"
	"strings"
)

func Parse(text string) (IRange, error) {
	if index := strings.IndexByte(text, '/'); index != -1 {
		if _, network, err := net.ParseCIDR(text); err == nil {
			return IpNetWrapper{network}, nil
		} else {
			return nil, err
		}
	}
	if ip := ParseIp(text); ip != nil {
		return IpWrapper{ip}, nil
	}
	if index := strings.IndexByte(text, '-'); index != -1 {
		if start, end := ParseIp(text[:index]), ParseIp(text[index+1:]); start != nil && end != nil {
			if len(start) == len(end) && !lessThan(end, start) {
				return &Range{start: start, end: end}, nil
			}
		}
		return nil, &net.ParseError{Type: "range", Text: text}
	}
	return nil, &net.ParseError{Type: "ip/CIDR address/range", Text: text}
}

func ParseIp(str string) net.IP {
	for _, b := range str {
		switch b {
		case '.':
			return net.ParseIP(str).To4()
		case ':':
			return net.ParseIP(str).To16()
		}
	}
	return nil
}
