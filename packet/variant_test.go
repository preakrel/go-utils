package packet

import (
	"errors"
	"testing"
)

func errJudge(err error) bool {
	return errors.Is(err, PackFormatError) || errors.Is(err, PackNilError) || errors.Is(err, URLParseError)
}

func TestVariant(t *testing.T) {
	type args struct {
		status VariantType
		p      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "URL2Body",
			args: args{
				status: URL2Body,
				p:      "",
			},
			want: "",
		},

		{
			name: "URL2Cookie",
			args: args{
				status: URL2Cookie,
				p:      "",
			},
			want: "",
		},

		{
			name: "Body2Cookie",
			args: args{
				status: Body2Cookie,
				p:      "",
			},
			want: "",
		},

		{
			name: "Body2URL",
			args: args{
				status: Body2URL,
				p:      "",
			},
			want: "",
		},

		{
			name: "Cookie2URL",
			args: args{
				status: Cookie2URL,
				p:      "",
			},
			want: "",
		},

		{
			name: "Cookie2Body",
			args: args{
				status: Cookie2Body,
				p:      "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Variant(uint8(tt.args.status), tt.args.p)
			if errJudge(err) {
				t.Log(err)
				return
			}
			if got != tt.want {
				t.Errorf("Variant() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantURL2Body(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "GET http://www.baidu.com/?x=1 HTTP/1.1\nAccept-Language: en\n\n"},
			want: "POST http://www.baidu.com/ HTTP/1.1\nAccept-Language: en\nContent-Type: application/x-www-form-urlencoded\n\nx=1",
		},
		{
			name: "success1",
			args: args{pkt: "GET /?x=1 HTTP/1.1\r\nAccept-Language: en\r\nHost: www.baidu.com\r\n\r\n"},
			want: "POST / HTTP/1.1\r\nAccept-Language: en\r\nHost: www.baidu.com\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\nx=1",
		},
		{
			name: "errLineBreak",
			args: args{pkt: "GET http://www.baidu.com/?x=1 HTTP/1.1\nAccept-Language: en\n"},
			want: "GET http://www.baidu.com/?x=1 HTTP/1.1\nAccept-Language: en\n",
		},
		{
			name: "errFirst",
			args: args{pkt: "GET\nAccept-Language: en\n\n"},
			want: "GET\nAccept-Language: en\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLToBody(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("URLToBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("URLToBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantURL2Cookie(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "GET http://www.baidu.com/?x=1 HTTP/1.1\nAccept-Language: en\n\n"},
			want: "GET http://www.baidu.com/ HTTP/1.1\nAccept-Language: en\nCookie: x=1\n\n",
		},

		{
			name: "success1",
			args: args{pkt: "GET http://www.baidu.com/?x=1&y=2 HTTP/1.1\nAccept-Language: en\nCookie: x=1\n\n"},
			want: "GET http://www.baidu.com/ HTTP/1.1\nAccept-Language: en\nCookie: x=1; x=1; y=2\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLToCookie(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("URLToCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLToCookie() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantBody2Cookie(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "POST / HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/x-www-form-urlencoded\n\nc=3&x=1&y=2"},
			want: "POST / HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/x-www-form-urlencoded\nCookie: c=3; x=1; y=2\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BodyToCookie(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("BodyToCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BodyToCookie() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantBody2URL(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "POST /?&y=2 HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/x-www-form-urlencoded\n\nc=3&x=1"},
			want: "POST /?c=3&x=1&y=2 HTTP/1.1\nHost: www.baidu.com\nContent-Type: application/x-www-form-urlencoded\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BodyToURL(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("BodyToURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BodyToURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantCookie2URL(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "GET http://www.baidu.com/?x=1&y=2 HTTP/1.1\nAccept-Language: en\nCookie: c=1\n\n"},
			want: "GET http://www.baidu.com/?c=1&x=1&y=2 HTTP/1.1\nAccept-Language: en\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CookieToURL(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("CookieToURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CookieToURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantCookie2Body(t *testing.T) {
	type args struct {
		pkt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{pkt: "GET http://www.baidu.com/?x=1&y=2 HTTP/1.1\nAccept-Language: en\nCookie: c=3; x=1; y=2\n\n"},
			want: "GET http://www.baidu.com/?x=1&y=2 HTTP/1.1\nAccept-Language: en\n\nc=3&x=1&y=2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CookieToBody(tt.args.pkt)
			if errJudge(err) {
				t.Log(err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("CookieToBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CookieToBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
