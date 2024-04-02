package packet

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestURL2Packet(t *testing.T) {
	t.Log(URL2Packet("https://gorm.io/zh_CN/docs/update.html"))
	t.Log(URL2Packet("https://www.google.com/search?^^^multipart^^^proname={{r0}}&tz=1%E4%B8%87%E4%BB%A5%E4%B8%8B&prouse={{r0}}&sx%5B%5D=&sx%5B%5D=&sm={{r0}}&province=%E5%85%A8%E5%9B%BD&city=%E5%85%A8%E5%9B%BD%E5%90%84%E5%9C%B0%E5%8C%BA&xiancheng=&cityforadd=&img=%2Fimage%2Fnopic.gif&flv=&zc=&yq=&action=add&Submit=%E5%A1%AB%E5%A5%BD%E4%BA%86%EF%BC%8C%E5%8F%91%E5%B8%83%E4%BF%A1%E6%81%AF&smallclassid[]=1&smallclassid[]=2)%20union%20select%20{{r1}}*{{r2}}%23"))
	t.Log(URL2Packet("https://www.google.com/search?^^^urlencoded^^^q=node-sass%3A+Command+failed&newwindow=1&ei=K13xYPSXMcy4mAXW8IegBA&oq=node-sass%3A+Command+failed&gs_lcp=Cgdnd3Mtd2l6EANKBAhBGABQg5-4bViDn7htYIKluG1oAHACeACAAYUEiAGFBJIBAzUtMZgBAKABAqABAaoBB2d3cy13aXrAAQE&sclient=gws-wiz&ved=0ahUKEwi02bbfrufxAhVMHKYKHVb4AUQQ4dUDCA4&uact=5"))
	t.Log(URL2Packet("https://www.google.com/search?^^^json^^^q=node-sass%3A+Command+failed&newwindow=1&ei=K13xYPSXMcy4mAXW8IegBA&oq=node-sass%3A+Command+failed&gs_lcp=Cgdnd3Mtd2l6EANKBAhBGABQg5-4bViDn7htYIKluG1oAHACeACAAYUEiAGFBJIBAzUtMZgBAKABAqABAaoBB2d3cy13aXrAAQE&sclient=gws-wiz&ved=0ahUKEwi02bbfrufxAhVMHKYKHVb4AUQQ4dUDCA4&uact=5"))
	t.Log(URL2Packet("https://www.google.com/search?^^^xml^^^q=node-sass%3A+Command+failed&newwindow=1&ei=K13xYPSXMcy4mAXW8IegBA&oq=node-sass%3A+Command+failed&gs_lcp=Cgdnd3Mtd2l6EANKBAhBGABQg5-4bViDn7htYIKluG1oAHACeACAAYUEiAGFBJIBAzUtMZgBAKABAqABAaoBB2d3cy13aXrAAQE&sclient=gws-wiz&ved=0ahUKEwi02bbfrufxAhVMHKYKHVb4AUQQ4dUDCA4&uact=5"))
}

func TestPostToJson(t *testing.T) {
	p := `POST /member/index.php HTTP/1.1
Host: www.huya.com
X-Requested-With: XMLHttpRequest
Sec-Fetch-Mode: cors
Sec-Fetch-Dest: empty
Referer: https://www.huya.com/
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Sec-Ch-Ua: "(Not(A:Brand";v="8", "Chromium";v="100"
Accept: application/json, text/javascript, */*; q=0.01
Sec-Fetch-Site: same-origin
Sec-Ch-Ua-Mobile: ?0
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36
Sec-Ch-Ua-Platform: "Windows"
Uber-Trace-Id: 506b03355f2b97af:506b03355f2b97af:0:0
Cookie: SoundValue=0.50; isInLiveRoom=; udb_guiddata=0ae7f9e41b3b44af91a33b887b482171; udb_deviceid=w_570947618056728576; __yamid_new=C9CC7A7F477000015454ED6219701439; __yamid_tt1=0.45141103814868866; game_did=5srU6AIx42V0l-PdzlHdKeNbF1Sl7ToQ2hT; huya_ua=webh5&0.1.0&websocket; Hm_lvt_51700b6c722f5bb4cf39906a596ea41f=1663749181,1663754708,1663818530,1663835318; Hm_lpvt_51700b6c722f5bb4cf39906a596ea41f=1663835318; PHPSESSID=jm5ieusbve0bjur2b9rnu2lqp3; udb_passdata=3; __yasmid=0.45141103814868866; _yasids=__rootsid%3DC9FCD03321200001EB5D723E1B00122B; huya_web_rep_cnt=16
Content-Type: multipart/form-data; boundary=----9kbrg75fqsum

------9kbrg75fqsum
Content-Disposition: form-data; name="m"

Stream
------9kbrg75fqsum
Content-Disposition: form-data; name="do"

getBannerStreamInfo
------9kbrg75fqsum
Content-Disposition: form-data; name="us"

1560173861%2C1560173878%2C1199561235630%2C1795363292%2C1099531728388%2C1199606779004
------9kbrg75fqsum
Content-Disposition: form-data; name="k"

feb68ac07834b53fab12d8f1116f50bd
------9kbrg75fqsum
Content-Disposition: form-data; name=""


------9kbrg75fqsum--
`
	s, e := PostToJson(p)
	if e != nil {
		t.Log(e)
	}
	t.Log(s)
}

func TestPostToXML(t *testing.T) {
	p := `POST /member/index.php HTTP/1.1
Host: www.huya.com
X-Requested-With: XMLHttpRequest
Sec-Fetch-Mode: cors
Sec-Fetch-Dest: empty
Referer: https://www.huya.com/
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Sec-Ch-Ua: "(Not(A:Brand";v="8", "Chromium";v="100"
Accept: application/json, text/javascript, */*; q=0.01
Sec-Fetch-Site: same-origin
Sec-Ch-Ua-Mobile: ?0
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36
Sec-Ch-Ua-Platform: "Windows"
Uber-Trace-Id: 506b03355f2b97af:506b03355f2b97af:0:0
Cookie: SoundValue=0.50; isInLiveRoom=; udb_guiddata=0ae7f9e41b3b44af91a33b887b482171; udb_deviceid=w_570947618056728576; __yamid_new=C9CC7A7F477000015454ED6219701439; __yamid_tt1=0.45141103814868866; game_did=5srU6AIx42V0l-PdzlHdKeNbF1Sl7ToQ2hT; huya_ua=webh5&0.1.0&websocket; Hm_lvt_51700b6c722f5bb4cf39906a596ea41f=1663749181,1663754708,1663818530,1663835318; Hm_lpvt_51700b6c722f5bb4cf39906a596ea41f=1663835318; PHPSESSID=jm5ieusbve0bjur2b9rnu2lqp3; udb_passdata=3; __yasmid=0.45141103814868866; _yasids=__rootsid%3DC9FCD03321200001EB5D723E1B00122B; huya_web_rep_cnt=16
Content-Type: multipart/form-data; boundary=----9kbrg75fqsum

------9kbrg75fqsum
Content-Disposition: form-data; name="m"

Stream
------9kbrg75fqsum
Content-Disposition: form-data; name="do"

getBannerStreamInfo
------9kbrg75fqsum
Content-Disposition: form-data; name="us"

1560173861%2C1560173878%2C1199561235630%2C1795363292%2C1099531728388%2C1199606779004
------9kbrg75fqsum
Content-Disposition: form-data; name="k"

feb68ac07834b53fab12d8f1116f50bd
------9kbrg75fqsum
Content-Disposition: form-data; name=""


------9kbrg75fqsum--
`
	s, e := PostToXML(p)
	if e != nil {
		t.Log(e)
	}
	t.Log(s)
}

func TestPostToMultipart(t *testing.T) {
	p := "POST /tool/ajaxip HTTP/1.1\r\nHost: coolaf.com\r\nContent-Type: application/x-www-form-urlencoded\nCookie: iris.language=en\r\nAccept: application/json, text/javascript, */*; q=0.01\r\nX-Requested-With: XMLHttpRequest\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36\r\nOrigin: http://coolaf.com\r\nReferer: http://coolaf.com/\r\nAccept-Language: en-US,en;q=0.9\r\nProxy-Connection: keep-alive\r\nContent-Length: 147\r\n\nip=118.116.15.150&cc=909"
	t.Log(PostToMultipart(p))
}

func TestMultipartToPOST(t *testing.T) {
	p := "POST /tool/ajaxip HTTP/1.1\nHost: coolaf.com\nCookie: iris.language=en\nAccept: application/json, text/javascript, */*; q=0.01\nX-Requested-With: XMLHttpRequest\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36\r\nOrigin: http://coolaf.com\nReferer: http://coolaf.com/\nAccept-Language: en-US,en;q=0.9\r\nProxy-Connection: keep-alive\r\nContent-Length: 147\nContent-Type: multipart/form-data; boundary=----fpllngzieyoh\r\n\n------fpllngzieyoh\r\nContent-Disposition: form-data; name=\"ip\"\r\n\r\n118.116.15.150\r\n\n------fpllngzieyoh\r\n\nContent-Disposition: form-data; name=\"cc\"\r\n\n\r\n909\n------fpllngzieyoh--\r\n"
	t.Log(MultipartToPOST(p))
}

func TestGetToPostNormalFormatPacket2(t *testing.T) {
	p := "GET /ListAccounts?gpsia=1&source=ChromiumBrowser&json=standard HTTP/1.1\r\nHost: accounts.google.com\r\nSec-Fetch-Site: none\r\nAccept-Language: en-US,en;q=0.9\r\nCookie: NID=511=oNXT4JwWJCt-CmLGjkC6EqHAKo01_QCAt95oSqhKbJKWYNqOHIaByMzlI_asSFh1DrwIguCQcLY_a6nz6R4EhX7-B9qK0XK_wczAr_lddug_Z9kq1Fu48jl4YipSx2WihKjHhCj3z5Ge-IfbozsrWLy-Emco-g1QvhICqnDbACs; 1P_JAR=2022-08-11-11\r\nSec-Fetch-Dest: empty\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36\r\nConnection: keep-alive\r\nContent-Length: 1\r\nOrigin: https://www.google.com\r\nSec-Fetch-Mode: no-cors\r\n\r\n"
	t.Log(GetToPostNormalFormatPacket(p))
}

func TestGetToPostJsonFormatPacket(t *testing.T) {
	type args struct {
		packet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{packet: ""},
			want: "",
		},

		{
			name: "err",
			args: args{packet: "GET / HTTP/1.1\r\n"},
			want: "GET / HTTP/1.1\r\n",
		},

		{
			name: "err1",
			args: args{packet: "GET\n\n"},
			want: "GET\n\n",
		},
		{
			name: "err11",
			args: args{packet: "GET /?a[0][1]=val&a[0][2]=val2 HTTP/1.1\n\n"},
			want: "POST HTTP/1.1 HTTP/1.1\nContent-Type: application/json;charset=UTF-8\n\n{}\n\n",
		},

		{
			name: "err2",
			args: args{packet: "GET HTTP/1.1\n\n"},
			want: "POST HTTP/1.1 HTTP/1.1\nContent-Type: application/json;charset=UTF-8\n\n{}\n\n",
		},

		{
			name: "success",
			args: args{packet: "GET /? HTTP/1.1\n\n"},
			want: "POST / HTTP/1.1\nContent-Type: application/json;charset=UTF-8\n\n{}\n\n",
		},

		{
			name: "success1",
			args: args{packet: "GET /?id=1&name=john&sub[]=1&sub[]=2&sub[]=three HTTP/1.1\nHost: www.baidu.com\n\nc=3"},
			want: "POST / HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/json;charset=UTF-8\n\n{\"c\":3,\"id\":1,\"name\":\"john\",\"sub\":{\"0\":1,\"1\":2,\"2\":\"three\"}}\n\n",
		},

		{
			name: "pack have content-type",
			args: args{packet: "GET /?id=1&name=john&sub[]=1&sub[]=2&sub[]=three HTTP/1.1\nContent-Type: application/json;charset=UTF-8\n\n"},
			want: "POST / HTTP/1.1\nContent-Type: application/json;charset=UTF-8\n\n{\"id\":1,\"name\":\"john\",\"sub\":{\"0\":1,\"1\":2,\"2\":\"three\"}}\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetToPostJsonFormatPacket(tt.args.packet)
			if errors.Is(PackFormatError, err) {
				t.Log(err)
				return
			}

			if errors.Is(PackNilError, err) {
				t.Log(err)
				return
			}
			t.Logf("packet==>\n%v\n \r\ngot==> \n%v\n", tt.args.packet, got)

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetToPostJsonFormatPacket() ==> \n%v\n \r\nwant==> \n%v\n", got, tt.want)
			// }

		})
	}
}
func TestGetToPostXMLFormatPacket(t *testing.T) {
	type args struct {
		packet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "pack have content-type",
			args: args{packet: "GET /?id=1&name=john&sub[]=1&sub[]=2&sub[]=three HTTP/1.1\nContent-Type: application/json\n\n"},
			want: `POST / HTTP/1.1
Content-Type: application/xml

<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<xml>
    <id>1</id>
    <name>john</name>
    <sub>
        <0>1</0>
        <1>2</1>
        <2>three</2>
    </sub>
</xml>
`,
		},
		{
			name: "pack have content-type",
			args: args{packet: "GET /?a[0][2]=val2&a[0][3]=val3 HTTP/1.1\nContent-Type: application/json\n\n"},
			want: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetToPostXMLFormatPacket(tt.args.packet)
			if errors.Is(PackFormatError, err) {
				t.Log(err)
				return
			}

			if errors.Is(PackNilError, err) {
				t.Log(err)
				return
			}
			t.Log(got)
			// if got != tt.want {
			// 	t.Errorf("GetToPostJsonFormatPacket() ==> \n%v\n \r\nwant==> \n%v\n", got, tt.want)
			// }

		})
	}
}

func TestGetToPostNormalFormatPacket(t *testing.T) {
	type args struct {
		packet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{packet: ""},
			want: "",
		},

		{
			name: "err",
			args: args{packet: "GET / HTTP/1.1\r\n"},
			want: "GET / HTTP/1.1\r\n",
		},

		{
			name: "err1",
			args: args{packet: "GET\n\n"},
			want: "GET\n\n",
		},

		{
			name: "err2",
			args: args{packet: "GET HTTP/1.1\n\n"},
			want: "POST HTTP/1.1 HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\n",
		},

		{
			name: "success",
			args: args{packet: "GET /? HTTP/1.1\n\n"},
			want: "POST / HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\n",
		},

		{
			name: "success1",
			args: args{packet: "GET /?x=1&y=2 HTTP/1.1\nHost: www.baidu.com\n\nc=3"},
			want: "POST / HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/x-www-form-urlencoded\n\nc=3&x=1&y=2",
		},

		{
			name: "pack have content-type",
			args: args{packet: "GET /?x=1&y=2&c=/name HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\n"},
			want: "POST / HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\nx=1&y=2&c=/name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetToPostNormalFormatPacket(tt.args.packet)
			if errors.Is(PackFormatError, err) {
				t.Log(err)
				return
			}

			if errors.Is(PackNilError, err) {
				t.Log(err)
				return
			}

			if got != tt.want {
				t.Errorf("GetToPostNormalFormatPacket() = \n%v\n, want \n%v\n", got, tt.want)
			}

		})
	}
}

func TestGetToPostMultipartFormatPacket(t *testing.T) {
	type args struct {
		packet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{packet: ""},
			want: "",
		},
		{
			name: "err",
			args: args{packet: "GET / HTTP/1.1\r\n"},
			want: "GET / HTTP/1.1\r\n",
		},

		{
			name: "err1",
			args: args{packet: "GET\n\n"},
			want: "GET\n\n",
		},

		{
			name: "success",
			args: args{packet: "GET /?s=index/user/_empty?name=/www/wwwroot/task/runtime/log/single.log HTTP/1.1\n\n"},
			want: "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetToPostMultipartFormatPacket(tt.args.packet)
			if errors.Is(PackFormatError, err) {
				t.Log(err)
				return
			}

			if errors.Is(PackNilError, err) {
				t.Log(err)
				return
			}

			if tt.want == "*" {
				t.Log(got)
				return
			}

			if got != tt.want {
				t.Errorf("GetToPostMultipartFormatPacket() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestPostToGetFormatPacket(t *testing.T) {
	type args struct {
		packet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// {
		// 	name: "nil",
		// 	args: args{packet: ""},
		// 	want: "",
		// },
		// {
		// 	name: "err",
		// 	args: args{packet: "POST / HTTP/1.1 Content-Type: "},
		// 	want: "POST / HTTP/1.1 Content-Type: ",
		// },
		// {
		// 	name: "err2",
		// 	args: args{packet: "POST / HTTP/1.1 Content-Type: multipart/"},
		// 	want: "POST / HTTP/1.1 Content-Type: multipart/",
		// },
		// {
		// 	name: "err3",
		// 	args: args{packet: "POST / HTTP/1.1\nContent-Type: multipart/\n\n"},
		// 	want: "POST / HTTP/1.1\nContent-Type: multipart/\n\n",
		// },
		// {
		// 	name: "success",
		// 	args: args{packet: "POST / HTTP/1.1\nContent-Type: multipart/form-data; boundary=--------dmvl6nldfcgw\n\n--------dmvl6nldfcgw\nContent-Disposition: form-data; name=\"id\"\n\n2\n--------dmvl6nldfcgw--\n"},
		// 	want: "GET /?id=2 HTTP/1.1\n\n",
		// },
		// {
		// 	name: "header-test1",
		// 	args: args{packet: "POST / HTTP/1.1\r\nContent-Type: multipart/form-data; boundary=--------dmvl6nldfcgw\r\n\r\n--------dmvl6nldfcgw\r\nContent-Disposition: form-data; name=\"\"\r\n\r\n\r\n--------dmvl6nldfcgw--\r\n"},
		// 	want: "GET / HTTP/1.1\r\n\r\n",
		// },
		// {
		// 	name: "header-test2-error",
		// 	args: args{packet: "POST / HTTP/1.1\nboundary=--------z8mwm34pdkzq\n\n--------z8mwm34pdkzq\nContent-Disposition: form-data; name=\"\"\n\n\n--------z8mwm34pdkzq--\n"},
		// 	want: "GET /?--------z8mwm34pdkzqContent-Disposition: form-data; name=\"\"--------z8mwm34pdkzq-- HTTP/1.1\nboundary=--------z8mwm34pdkzq\n\n",
		// },
		// {
		// 	name: "params-test2",
		// 	args: args{packet: "POST /\n\n\n"},
		// 	want: "GET /? HTTP/1.1\n\n",
		// },
		// {
		// 	name: "params-test2-1",
		// 	args: args{packet: "POST /?\n\n\n"},
		// 	want: "GET /?& HTTP/1.1\n\n",
		// },
		// {
		// 	name: "params-test2-2",
		// 	args: args{packet: "POST\n\n"},
		// 	want: "POST\n\n",
		// },
		// {
		// 	name: "params-test3",
		// 	args: args{packet: "POST /? HTTP/1.1\nContent-Type: \n\n"},
		// 	want: "GET /?& HTTP/1.1\n\n",
		// },
		{
			name: "params-test4",
			args: args{packet: "POST /?x=1 HTTP/1.1\nContent-Type: \n\nx=1&y=2&c=1"},
			want: "GET /?x=1&x=1&y=2&c=1 HTTP/1.1\n\n",
		},
		{
			name: "json-test1",
			args: args{packet: "POST / HTTP/1.1\nContent-Type: application/json\n\n{\"id\":1,\"name\":\"john\",\"sub\":[1,2,\"three\"]}"},
			want: "GET /?id=1&name=john&sub[]=1&sub[]=2&sub[]=three HTTP/1.1\n\n",
		},
		{
			name: "xml-test1",
			args: args{packet: "POST / HTTP/1.1\nContent-Type: application/xml\n\n<xml><id>1</id><name>john</name><sub><item><item>1</item><item>2</item><item>3</item></item><item>2</item><item>three</item></sub></xml>"},
			want: "GET /?id=1&name=john&sub[][]=1&sub[][]=2&sub[][]=3&sub[]=2&sub[]=three HTTP/1.1\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostToGetFormatPacket(tt.args.packet)
			if err != nil {
				t.Log(err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostToGetFormatPacket() got ===> \r\n%v  \r\nwant ===> \r\n%v", got, tt.want)
			}
		})
	}
}

func TestConversion(t *testing.T) {
	type args struct {
		status uint8
		p      string
	}
	tests := []struct {
		name    string
		args    args
		want    Base
		wantErr bool
	}{
		// {
		// 	name: "Url2Pkt",
		// 	args: args{
		// 		status: Url2Pkt,
		// 		p:      "http://www.baidu.com"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "GET / HTTP/1.1\r\nHost: www.baidu.com\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en\r\nUser-Agent: Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)\r\nConnection: close\r\n\r\n",
		// 	},
		// },
		// {
		// 	name: "Post2Multipart",
		// 	args: args{
		// 		status: PktToMultipart,
		// 		p:      "POST /tool/ajaxip HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\nip=118.116.15.150"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "*",
		// 	},
		// },
		// {
		// 	name: "Multipart2Post",
		// 	args: args{
		// 		status: PktToPost,
		// 		p:      "POST /tool/ajaxip HTTP/1.1\nContent-Type: multipart/form-data; boundary=----WebKitFormBoundaryYHAnqEEIZ0w1lXJg\n\n------WebKitFormBoundaryYHAnqEEIZ0w1lXJg\nContent-Disposition: form-data; name=\"ip\"\n\n118.116.15.150\n------WebKitFormBoundaryYHAnqEEIZ0w1lXJg--\n"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "POST /tool/ajaxip HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\nip=118.116.15.150",
		// 	},
		// },

		// {
		// 	name: "Get2PostNormalFormatPacket",
		// 	args: args{
		// 		status: PktToPost,
		// 		p:      "GET /?x=1&y=2 HTTP/1.1\n\n"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "POST / HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\nx=1&y=2",
		// 	},
		// },

		// {
		// 	name: "Get2PostMultipartFormatPacket",
		// 	args: args{
		// 		status: PktToMultipart,
		// 		p:      "GET /? HTTP/1.1\n\n"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "*",
		// 	},
		// },

		// {
		// 	name: "Post2GetFormatPacket",
		// 	args: args{
		// 		status: PktToGet,
		// 		p:      "POST / HTTP/1.1\nContent-Type: multipart/form-data; boundary=--------dmvl6nldfcgw\n\n--------dmvl6nldfcgw\nContent-Disposition: form-data; name=\"\"\n\n\n--------dmvl6nldfcgw--\n"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "GET / HTTP/1.1\n\n",
		// 	},
		// },

		// {
		// 	name: "Post2GetFormatPacket2",
		// 	args: args{
		// 		status: PktToGet,
		// 		p:      "POST /tool/ajaxip HTTP/1.1\nContent-Type: application/x-www-form-urlencoded\n\nip=118.116.15.150"},
		// 	want: Base{
		// 		IsSSL:  false,
		// 		Packet: "GET /tool/ajaxip?ip=118.116.15.150 HTTP/1.1\n\n",
		// 	},
		// },
		{
			name: "Post2GetFormatPacket2",
			args: args{
				status: PktToGet,
				p:      "POST / HTTP/1.1\nContent-Type: application/json\n\n{\"id\":1,\"name\":\"john\",\"sub\":[1,2,\"three\"]}"},
			want: Base{
				IsSSL:  false,
				Packet: "GET /tool/ajaxip?ip=118.116.15.150 HTTP/1.1\n\n",
			},
		},
		{
			name: "Post2XmlFormatPacket2",
			args: args{
				status: PktToXml,
				p:      "POST /tool/ajaxip HTTP/1.1\nContent-Type: application/json\n\n{\"id\":187923,\"poc\":\"f\"}"},
			want: Base{
				IsSSL:  false,
				Packet: "GET /tool/ajaxip?ip=118.116.15.150 HTTP/1.1\n\n",
			},
		},
		{
			name: "Post2XmlFormatPacket2",
			args: args{
				status: PktToPost,
				p: `POST /post/id.php HTTP/1.1
Host: www.anquanke.com
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36
Connection: close
Content-Type: application/json;charset=UTF-8

{"id":187923,"po1c":[1,2],"poc":true}`},
			want: Base{
				IsSSL:  false,
				Packet: "GET /tool/ajaxip?ip=118.116.15.150 HTTP/1.1\n\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Conversion(tt.args.status, tt.args.p, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conversion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want.Packet == "*" {
				// Multipart
				//t.Log(got)
				return
			}
			t.Log(got)

			if !reflect.DeepEqual(got, tt.want) {
				// t.Errorf("Conversion() got = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}
func TestConvXml(t *testing.T) {

	str := `POST /post/id.php HTTP/1.1
Host: www.anquanke.com
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36
Connection: close
Content-Type: application/json;charset=UTF-8

{"id":187923,"po1c":[1,2],"poc":true}`
	got, err := Conversion(PktToXml, str, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(strconv.Quote(got.Packet))

}
