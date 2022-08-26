package ipnetblocks

import (
	"encoding/json"
	"fmt"
	"time"
)

// unmarshalString parses the JSON-encoded data and returns value as a string.
func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time.
type Time time.Time

var emptyTime Time

// UnmarshalJSON decodes time as IP Netblocks API does.
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}

	v, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as IP Netblocks API does.
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}

// AS is the related autonomous system's data.
type AS struct {
	// ASN is the autonomous system number.
	ASN int `json:"asn"`

	// Name is the autonomous system name.
	Name string `json:"name"`

	// Type is the autonomous system type. One of the following: "Cable/DSL/ISP", "Content", "Educational/Research",
	// "Enterprise", "Non-Profit", "Not Disclosed", "NSP", "Route Server". Empty when unknown.
	Type string `json:"type"`

	// Route is the autonomous system route.
	Route string `json:"route"`

	// Domain is the autonomous system Website's URL.
	Domain string `json:"domain"`
}

// Contact is the contact data.
type Contact struct {
	// ID is the contact's ID.
	ID string `json:"id"`

	// Person is the name of the contact person. Indicates that object is of person type.
	// Fields Person and Role are mutually exclusive.
	Person string `json:"person"`

	// Role is the name of the contact role.
	// Fields Person and Role are mutually exclusive.
	Role string `json:"role"`

	// Email is the email address.
	Email string `json:"email"`

	// Phone is the phone number.
	Phone string `json:"phone"`

	// Country is two letters country code from ISO 3166.
	Country string `json:"country"`

	// City is the name of city.
	City string `json:"city"`

	// Address is the location information. May not reflect exact physical location.
	Address []string `json:"address"`
}

// Organization is the organization data.
type Organization struct {
	// Org is the organisation's ID.
	Org string `json:"org"`

	// Name is the organisation name.
	Name string `json:"name"`

	// Phone is the phone number.
	Phone string `json:"phone"`

	// Email is the email address.
	Email string `json:"email"`

	// Country is two letters country code from ISO 3166.
	Country string `json:"country"`

	// City is the name of city.
	City string `json:"city"`

	// PostalCode is the postal code.
	PostalCode string `json:"postalCode"`

	// Address is the location information. May not reflect exact physical location.
	Address []string `json:"address"`
}

// Maintainer is the maintainer data.
type Maintainer struct {
	// Mntner is the maintainer's ID.
	Mntner string `json:"mntner"`

	// Email is the email address.
	Email string `json:"email"`
}

// Inetnum is a part of IP Netblocks API response that contains the netblock (inetnum) related data.
type Inetnum struct {
	// Inetnum is the IP range.
	Inetnum string `json:"inetnum"`

	// InetnumFirst is the first IP as 128-bit unsigned integer value, stored as floating-point number.
	InetnumFirst float64 `json:"inetnumFirst"`

	// InetnumLast the last IP as 128-bit unsigned integer value, stored as floating-point number.
	InetnumLast float64 `json:"inetnumLast"`

	// InetnumFirstString is the string representation of the inetnumFirst field.
	// Use this field if you want to avoid exponential representations of the inetnumFirst field.
	InetnumFirstString string `json:"inetnumFirstString"`

	// InetnumLastString is the string representation of the inetnumLast field.
	// Use this field if you want to avoid exponential representations of the inetnumLast field.
	InetnumLastString string `json:"inetnumLastString"`

	// Parent is the reference to the block from which the information was borrowed.
	// Its presence indicates that the block was obtained from BGP routing tables.
	Parent string `json:"parent"`

	// AS is the related autonomous system's data.
	AS AS `json:"as"`

	// Netname is the name of the IPs range.
	Netname string `json:"netname"`

	// Nethandle is the ID of the block from ARIN.
	Nethandle string `json:"nethandle"`

	// Description is the description related to the block.
	Description []string `json:"description"`

	// Modified is the time when the IP Netblock was modified the last time, accordingly to the information
	// provided by the registry.
	Modified Time `json:"modified"`

	// Country is two letters country code from ISO 3166.
	Country string `json:"country"`

	// City is the name of city.
	City string `json:"city"`

	// Address is the location information. May not reflect exact physical location.
	Address []string `json:"address"`

	// AbuseContact is the list of abuse contacts.
	AbuseContact []Contact `json:"abuseContact"`

	// AdminContact is the list of administrative contacts.
	AdminContact []Contact `json:"adminContact"`

	// TechContact is the list of technical contacts.
	TechContact []Contact `json:"techContact"`

	// Org is the organisation registered the range.
	Org Organization `json:"org"`

	// MntBy is the list of maintainers who are able to update the IPs range.
	MntBy []Maintainer `json:"mntBy"`

	// MntDomains is the list of domains' maintainers.
	MntDomains []Maintainer `json:"mntDomains"`

	// MntLower is the list of maintainers who are able to change sub ranges.
	MntLower []Maintainer `json:"mntLower"`

	// MntRoutes is the list of maintainers of routing info.
	MntRoutes []Maintainer `json:"mntRoutes"`

	// Remarks is remarks and comments associated with the IP Netblock.
	Remarks []string `json:"remarks"`

	// Source is the source of range.
	Source string `json:"source"`
}

// Result is a part of the IP Netblock API response.
type Result struct {
	// Count is the number of records returned.
	Count int `json:"count"`

	//Limit is the maximum number of elements.
	Limit int `json:"limit"`

	// From is the IP netblock range that was passed in the request as an offset for the returned results.
	From *string `json:"from"`

	// Next If not null, it shows the value of the last IP netblock (inetnum) from the response, which can be
	// substituted into the From parameter and get the next page, if the current response doesn't fit into the limit.
	// Such pagination is only available for IPv4, IPv6 and ASN API requests. ORG requests don't support it.
	Next *string `json:"next"`

	// Inetnums is the list of returned netblocks.
	Inetnums []Inetnum `json:"inetnums"`
}

// IPNetblocksResponse is a response of IP Netblocks API.
type IPNetblocksResponse struct {
	// Search is a normalized search term.
	Search string `json:"search"`

	// Result is a result. This field is omitted in case of error.
	Result Result `json:"result"`

	// Error is the error message. This field is omitted when a call is successful.
	Error string `json:"error"`
}

// ErrorMessage is an error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
