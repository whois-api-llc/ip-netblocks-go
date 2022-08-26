package ipnetblocks

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathIPNetblocksResponseOK         = "/IPNetblocks/ok"
	pathIPNetblocksResponseError      = "/IPNetblocks/error"
	pathIPNetblocksResponse500        = "/IPNetblocks/500"
	pathIPNetblocksResponsePartial1   = "/IPNetblocks/partial"
	pathIPNetblocksResponsePartial2   = "/IPNetblocks/partial2"
	pathIPNetblocksResponseUnparsable = "/IPNetblocks/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the IP Netblocks API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathIPNetblocksResponseOK:
		case pathIPNetblocksResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathIPNetblocksResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathIPNetblocksResponsePartial1:
			response = response[:len(response)-10]
		case pathIPNetblocksResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathIPNetblocksResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new IP Netblocks API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:         apiServer.Client(),
		IPNetblocksBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestIPNetblocksGetByIP tests the GetByIP function.
func TestIPNetblocksGetByIP(t *testing.T) {
	checkResultRec := func(res *IPNetblocksResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IP
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: `API error: [499] Test error message.`,
		},
		{
			name: "invalid argument2",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{},
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "ip" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.GetByIP(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetByIP() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("IPNetblocks.GetByIP() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("IPNetblocks.GetByIP() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestIPNetblocksGetByCIDR tests the GetByIP function.
func TestIPNetblocksGetByCIDR(t *testing.T) {
	checkResultRec := func(res *IPNetblocksResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IPNet
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(100),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{}, net.IPMask{}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: `invalid argument: "ip" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.GetByCIDR(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetByCIDR() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("IPNetblocks.GetByCIDR() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("IPNetblocks.GetByCIDR() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestIPNetblocksGetByASN tests the GetByASN function.
func TestIPNetblocksGetByASN(t *testing.T) {
	checkResultRec := func(res *IPNetblocksResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory int
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					9876543210,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "9876543210" is invalid autonomous system number`,
		},
		{
			name: "invalid argument2",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					-1,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "-1" is invalid autonomous system number`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.GetByASN(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetByASN() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("IPNetblocks.GetByASN() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("IPNetblocks.GetByASN() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestIPNetblocksGetByOrg tests the GetByOrg function.
func TestIPNetblocksGetByOrg(t *testing.T) {
	checkResultRec := func(res *IPNetblocksResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "org" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.GetByOrg(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetByOrg() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("IPNetblocks.GetByOrg() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("IPNetblocks.GetByOrg() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestIPNetblocksGetRawByIP tests the GetRawByIP function.
func TestIPNetblocksGetRawByIP(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IP
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{8, 8, 8},
					OptionLimit(100),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument2",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IP{},
					OptionLimit(100),
				},
			},
			wantErr: "invalid argument: \"ip\" can not be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRawByIP(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetRawByIP() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("IPNetblocks.GetRawByIP() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}

// TestIPNetblocksGetRawByCIDR tests the GetByIP function.
func TestIPNetblocksGetRawByCIDR(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory net.IPNet
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(100),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(100),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 499",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{8, 8, 8, 8}, net.IPMask{0xFF, 0xFF, 0xFF, 0x0}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					net.IPNet{net.IP{}, net.IPMask{}},
					OptionLimit(1),
				},
			},
			want:    false,
			wantErr: `invalid argument: "ip" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRawByCIDR(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetRawByCIDR() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("IPNetblocks.GetRawByCIDR() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}

// TestIPNetblocksGetRawByASN tests the GetRawByASN function.
func TestIPNetblocksGetRawByASN(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory int
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 499",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					1234,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					9876543210,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "9876543210" is invalid autonomous system number`,
		},
		{
			name: "invalid argument2",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					-1,
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "-1" is invalid autonomous system number`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRawByASN(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetRawByASN() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("IPNetblocks.GetRawByASN() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}

// TestIPNetblocksGetRawByOrg tests the GetByOrg function.
func TestIPNetblocksGetRawByOrg(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"search":"8.8.8.8","result":{"count":1,"limit":1,"from":"","next":"8.8.8.0-8.8.8.255",
"inetnums":[{"inetnum":"8.8.8.0 - 8.8.8.255","inetnumFirst":281470816487424,"inetnumLast":281470816487679,
"inetnumFirstString":"281470816487424","inetnumLastString":"281470816487679","as":{"asn":15169,"name":"GOOGLE",
"type":"Content","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/"},"netname":"LVLT-GOGL-8-8-8",
"nethandle":"NET-8-8-8-0-1","description":[],"modified":"2014-03-14T00:00:00Z","country":"US","city":"Mountain View",
"address":["1600 Amphitheatre Parkway"],"abuseContact":[],"adminContact":[],"techContact":[],"org":{"org":"GOGL",
"name":"Google LLC","email":"arin-contact@google.com\nnetwork-abuse@google.com","phone":"+1-650-253-0000",
"country":"US","city":"Mountain View","postalCode":"94043","address":["1600 Amphitheatre Parkway"]},"mntBy":[],
"mntDomains":[],"mntLower":[],"mntRoutes":[],"remarks":[],"source":"ARIN"}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathIPNetblocksResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathIPNetblocksResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathIPNetblocksResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathIPNetblocksResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "API failed with status code: 499",
		},
		{
			name: "unparsable response",
			path: pathIPNetblocksResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"Whoisxmlapi",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: "",
		},
		{
			name: "invalid argument1",
			path: pathIPNetblocksResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionLimit(10),
				},
			},
			want:    false,
			wantErr: `invalid argument: "org" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRawByOrg(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("IPNetblocks.GetRawByOrg() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("IPNetblocks.GetRawByOrg() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
