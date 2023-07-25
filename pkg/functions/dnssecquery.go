package functions

import (
	"fmt"

	"github.com/miekg/dns"
)

func GetRRSIGWithKey(domain string) ([]*dns.RRSIG, []*dns.DNSKEY, error) {
	// Create a DNS client
	client := dns.Client{}

	// Prepare the DNS message for RRSIG query
	rrsigRequest := dns.Msg{}
	rrsigRequest.SetQuestion(domain+".", dns.TypeRRSIG)
	rrsigResponse, _, err := client.Exchange(&rrsigRequest, "8.8.8.8:53")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query RRSIG: %v", err)
	}

	// Prepare the DNS message for DNSKEY query
	dnskeyRequest := dns.Msg{}
	dnskeyRequest.SetQuestion(domain+".", dns.TypeDNSKEY)
	dnskeyResponse, _, err := client.Exchange(&dnskeyRequest, "8.8.8.8:53")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query DNSKEY: %v", err)
	}

	// Retrieve the RRSIG records
	var rrsigRecords []*dns.RRSIG
	for _, answer := range rrsigResponse.Answer {
		if rrsig, ok := answer.(*dns.RRSIG); ok {
			rrsigRecords = append(rrsigRecords, rrsig)
		}
	}

	// Retrieve the DNSKEY records
	var dnskeyRecords []*dns.DNSKEY
	for _, answer := range dnskeyResponse.Answer {
		if dnskey, ok := answer.(*dns.DNSKEY); ok {
			dnskeyRecords = append(dnskeyRecords, dnskey)
		}
	}

	return rrsigRecords, dnskeyRecords, nil
}
