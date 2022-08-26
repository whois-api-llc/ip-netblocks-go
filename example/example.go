package example

import (
	"context"
	"errors"
	ipnetblocks "github.com/whois-api-llc/ip-netblocks-go"
	"log"
	"net"
	"time"
)

func GetData(apikey string) {
	client := ipnetblocks.NewBasicClient(apikey)

	// Get parsed IP Netblocks API response by IP address as a model instance.
	ipNetblocksResp, resp, err := client.GetByIP(context.Background(),
		net.ParseIP("8.8.8.8"),
		// this option is ignored, as the inner parser works with JSON only.
		ipnetblocks.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *ipnetblocks.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print some values from each returned netblock.
	for _, obj := range ipNetblocksResp.Result.Inetnums {
		log.Printf("Netblock: %s, Time: %s, ASN: %d\n",
			obj.Inetnum,
			time.Time(obj.Modified).Format(time.RFC3339),
			obj.AS.ASN,
		)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetDataByASN(apikey string) {
	client := ipnetblocks.NewBasicClient(apikey)

	// Get parsed IP Netblocks API response by autonomous system number.
	ipNetblocksResp, _, err := client.GetByASN(context.Background(),
		15169)

	if err != nil {
		// Handle error message returned by server.
		var apiErr *ipnetblocks.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Just print the number of netblocks returned.
	log.Println(ipNetblocksResp.Result.Count)
}

func GetDataByOrg(apikey string) {
	client := ipnetblocks.NewBasicClient(apikey)

	// Get parsed IP Netblocks API response with IP netblocks which have the specified search term in their
	// Netblock (netname, description, remarks), or Organisation (org.org, org.name, org.email, org.address) fields.
	ipNetblocksResp, _, err := client.GetByOrg(context.Background(),
		"Amazon")

	if err != nil {
		// Handle error message returned by server.
		var apiErr *ipnetblocks.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print some values from each returned netblock.
	for _, obj := range ipNetblocksResp.Result.Inetnums {
		log.Printf("Netblock: %s, Organization: %s\n",
			obj.Inetnum,
			obj.Org.Name,
		)
	}
}

func GetAllDataByCIDR(apikey string) {
	client := ipnetblocks.NewBasicClient(apikey)

	var from *string
	var inetnums []ipnetblocks.Inetnum

	// Get implied network by CIDR notation.
	_, ipNet, err := net.ParseCIDR("8.8.8.8/10")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Get parsed IP Netblocks API response by CIDR.
		ipNetblocksResp, _, err := client.GetByCIDR(context.Background(), *ipNet,
			// this option defines maximum number of returned netblocks.
			ipnetblocks.OptionLimit(1000),
			// this option sets the IP netblock range that is used as an offset for the returned results.
			ipnetblocks.OptionFrom(from))
		if err != nil {
			// Handle error message returned by server.
			var apiErr *ipnetblocks.ErrorMessage
			if errors.As(err, &apiErr) {
				log.Println(apiErr.Code)
				log.Println(apiErr.Message)
			}
			log.Fatal(err)
		}

		// Store all returned netblocks in the single slice.
		inetnums = append(inetnums, ipNetblocksResp.Result.Inetnums...)

		// Break the loop when the last page is reached.
		if from = ipNetblocksResp.Result.Next; from == nil {
			break
		}
	}

	// Then print the count and some values from each netblock.
	log.Println(len(inetnums))
	for _, obj := range inetnums {
		log.Printf("Netblock: %s, Parent: %s, ASN: %d\n",
			obj.Inetnum,
			obj.Parent,
			obj.AS.ASN,
		)
	}
}

func GetRawDataByIP(apikey string) {
	client := ipnetblocks.NewBasicClient(apikey)

	// Get raw API response by IP for IPv6 address 2001:0000:4136:e378::.
	resp, err := client.GetRawByIP(context.Background(),
		net.ParseIP("2001:0000:4136:e378::"),
		ipnetblocks.OptionOutputFormat("JSON"))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
