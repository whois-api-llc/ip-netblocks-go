[![ip-netblocks-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![ip-netblocks-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/ip-netblocks-go)
[![ip-netblocks-go test](https://github.com/whois-api-llc/ip-netblocks-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/ip-netblocks-go/actions/)

# Overview

The client library for
[IP Netblocks API](https://ip-netblocks.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/ip-netblocks-go
```

# Examples

Full API documentation available [here](https://ip-netblocks.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := ipnetblocks.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := ipnetblocks.NewClient(apiKey, ipnetblocks.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

IP Netblocks API lets you get exhaustive information on the IP range that a given IP address belongs to.

```go

// Make request to get all parsed IP netblocks (inetnums) by IP address
ipNetblocksResp, resp, err := client.GetByIP(ctx, []byte{8,8,8,8})
if err != nil {
    log.Fatal(err)
}

for _, obj := range ipNetblocksResp.Result.Inetnums {
    log.Printf("Netblock: %s, Time: %s, ASN: %s\n",
        obj.Inetnum,
        time.Time(obj.Modified).Format(time.RFC3339),
        obj.AS.ASN,
    )
}

// Make request to get raw IP Netblocks data by autonomous system number
resp, err := client.GetRawByASN(context.Background(), 15169)
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
