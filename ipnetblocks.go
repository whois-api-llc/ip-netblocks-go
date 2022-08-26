package ipnetblocks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// IPNetblocks is an interface for IP Netblocks API.
type IPNetblocks interface {
	// GetByIP returns parsed IP Netblocks API response by IP address.
	GetByIP(ctx context.Context, ip net.IP, opts ...Option) (*IPNetblocksResponse, *Response, error)

	// GetByCIDR returns parsed IP Netblocks API response by CIDR.
	GetByCIDR(ctx context.Context, ip net.IPNet, opts ...Option) (*IPNetblocksResponse, *Response, error)

	// GetByASN returns parsed IP Netblocks API response by autonomous system number.
	GetByASN(ctx context.Context, asn int, opts ...Option) (*IPNetblocksResponse, *Response, error)

	// GetByOrg returns parsed IP Netblocks API response by organization.
	GetByOrg(ctx context.Context, org string, opts ...Option) (*IPNetblocksResponse, *Response, error)

	// GetRawByIP returns raw IP Netblocks API response by IP address as Response struct with Body saved
	// as a byte slice.
	GetRawByIP(ctx context.Context, ip net.IP, opts ...Option) (*Response, error)

	// GetRawByCIDR returns raw IP Netblocks API response by CIDR as Response struct with Body saved as a byte slice.
	GetRawByCIDR(ctx context.Context, ip net.IPNet, opts ...Option) (*Response, error)

	// GetRawByASN returns raw IP Netblocks API response by ASN as Response struct with Body saved as a byte slice.
	GetRawByASN(ctx context.Context, asn int, opts ...Option) (*Response, error)

	// GetRawByOrg returns raw IP Netblocks API response by organization as Response struct with Body saved
	// as a byte slice.
	GetRawByOrg(ctx context.Context, org string, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// ipNetblocksServiceOp is the type implementing the IPNetblocks interface.
type ipNetblocksServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ IPNetblocks = &ipNetblocksServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service ipNetblocksServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing IP Netblocks API response as a model instance.
type apiResponse struct {
	IPNetblocksResponse
	ErrorMessage
}

// request returns intermediate API response for further actions.
func (service ipNetblocksServiceOp) request(ctx context.Context, ip string, mask string, asn string, org string, opts ...Option) (*Response, error) {
	var err error

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if ip != "" {
		q.Set("ip", ip)
	}

	if mask != "" {
		q.Set("mask", mask)
	}

	if asn != "" {
		q.Set("asn", asn)
	}

	if org != "" {
		q.Set("org", org)
	}

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw IP Netblocks API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// validateASN validates autonomous system number.
func validateASN(asn int) (err error) {

	if asn < 0 || asn > 4294967295 {
		return &ArgError{fmt.Sprintf("%d", asn), "is invalid autonomous system number"}
	}

	return nil
}

// GetByIP returns parsed IP Netblocks API response by IP address.
func (service ipNetblocksServiceOp) GetByIP(
	ctx context.Context,
	ip net.IP,
	opts ...Option,
) (ipNetblocksResponse *IPNetblocksResponse, resp *Response, err error) {
	ipString := ip.String()
	if ipString == "<nil>" {
		return nil, nil, &ArgError{"ip", "can not be empty"}
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, ipString, "", "", "", optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	ipNetblocksResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if ipNetblocksResp.Message != "" || ipNetblocksResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    ipNetblocksResp.Code,
			Message: ipNetblocksResp.Message,
		}
	}

	return &ipNetblocksResp.IPNetblocksResponse, resp, nil
}

// GetByCIDR returns parsed IP Netblocks API response by CIDR.
func (service ipNetblocksServiceOp) GetByCIDR(
	ctx context.Context,
	ip net.IPNet,
	opts ...Option,
) (ipNetblocksResponse *IPNetblocksResponse, resp *Response, err error) {
	ipString := ip.IP.String()
	if ipString == "<nil>" {
		return nil, nil, &ArgError{"ip", "can not be empty"}
	}

	maskSize, _ := ip.Mask.Size()

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, ipString, fmt.Sprintf("%d", maskSize), "", "", optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	ipNetblocksResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if ipNetblocksResp.Message != "" || ipNetblocksResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    ipNetblocksResp.Code,
			Message: ipNetblocksResp.Message,
		}
	}

	return &ipNetblocksResp.IPNetblocksResponse, resp, nil
}

// GetByASN returns parsed IP Netblocks API response by autonomous system number.
func (service ipNetblocksServiceOp) GetByASN(
	ctx context.Context,
	asn int,
	opts ...Option,
) (ipNetblocksResponse *IPNetblocksResponse, resp *Response, err error) {
	if err = validateASN(asn); err != nil {
		return nil, nil, err
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, "", "", fmt.Sprintf("%d", asn), "", optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	ipNetblocksResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if ipNetblocksResp.Message != "" || ipNetblocksResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    ipNetblocksResp.Code,
			Message: ipNetblocksResp.Message,
		}
	}

	return &ipNetblocksResp.IPNetblocksResponse, resp, nil
}

// GetByOrg returns parsed IP Netblocks API response by organization.
func (service ipNetblocksServiceOp) GetByOrg(
	ctx context.Context,
	org string,
	opts ...Option,
) (ipNetblocksResponse *IPNetblocksResponse, resp *Response, err error) {
	if org == "" {
		return nil, nil, &ArgError{"org", "can not be empty"}
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, "", "", "", org, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	ipNetblocksResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if ipNetblocksResp.Message != "" || ipNetblocksResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    ipNetblocksResp.Code,
			Message: ipNetblocksResp.Message,
		}
	}

	return &ipNetblocksResp.IPNetblocksResponse, resp, nil
}

// GetRawByIP returns raw IP Netblocks API response by IP address as Response struct with Body saved
// as a byte slice.
func (service ipNetblocksServiceOp) GetRawByIP(
	ctx context.Context,
	ip net.IP,
	opts ...Option,
) (resp *Response, err error) {
	ipString := ip.String()
	if ipString == "<nil>" {
		return nil, &ArgError{"ip", "can not be empty"}
	}

	resp, err = service.request(ctx, ipString, "", "", "", opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// GetRawByCIDR returns raw IP Netblocks API response by CIDR as Response struct with Body saved as a byte slice.
func (service ipNetblocksServiceOp) GetRawByCIDR(
	ctx context.Context,
	ip net.IPNet,
	opts ...Option,
) (resp *Response, err error) {
	ipString := ip.IP.String()
	if ipString == "<nil>" {
		return nil, &ArgError{"ip", "can not be empty"}
	}

	maskSize, _ := ip.Mask.Size()

	resp, err = service.request(ctx, ipString, fmt.Sprintf("%d", maskSize), "", "", opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// GetRawByASN returns raw IP Netblocks API response by ASN as Response struct with Body saved as a byte slice.
func (service ipNetblocksServiceOp) GetRawByASN(
	ctx context.Context,
	asn int,
	opts ...Option,
) (resp *Response, err error) {
	if err = validateASN(asn); err != nil {
		return nil, err
	}

	resp, err = service.request(ctx, "", "", fmt.Sprintf("%d", asn), "", opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// GetRawByOrg returns raw IP Netblocks API response by organization as Response struct with Body saved
// as a byte slice.
func (service ipNetblocksServiceOp) GetRawByOrg(
	ctx context.Context,
	org string,
	opts ...Option,
) (resp *Response, err error) {
	if org == "" {
		return nil, &ArgError{"org", "can not be empty"}
	}

	resp, err = service.request(ctx, "", "", "", org, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
