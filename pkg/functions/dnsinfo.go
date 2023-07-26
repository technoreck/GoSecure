package functions

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type DNSResponse struct {
	A     []net.IP     `json:"A"`
	AAAA  []net.IP     `json:"AAAA"`
	MX    []*net.MX    `json:"MX"`
	TXT   []string     `json:"TXT"`
	NS    []string     `json:"NS"`
	CNAME []string     `json:"CNAME"`
	SOA   *SOAResponse `json:"SOA"`
	SRV   []*net.SRV   `json:"SRV"`
	PTR   []string     `json:"PTR"`
}

type SOAResponse struct {
	NS      string `json:"ns"`
	MBox    string `json:"mbox"`
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minttl  uint32 `json:"minttl"`
}

func PerformDNSLookup(hostname string) (DNSResponse, error) {
	resp := DNSResponse{}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return resp, err
	}
	resp.A = make([]net.IP, len(addrs))
	for i, addr := range addrs {
		resp.A[i] = net.ParseIP(addr)
	}

	aaaa, err := net.LookupIP(hostname)
	if err != nil {
		return resp, err
	}
	for _, addr := range aaaa {
		if addr.To4() == nil {
			resp.AAAA = append(resp.AAAA, addr)
		}
	}

	mx, err := net.LookupMX(hostname)
	if err == nil {
		resp.MX = mx
	} else {
		log.Println("MX lookup error:", err)
	}

	txt, err := net.LookupTXT(hostname)
	if err == nil {
		resp.TXT = txt
	} else {
		log.Println("TXT lookup error:", err)
	}

	ns, err := net.LookupNS(hostname)
	if err == nil {
		for _, n := range ns {
			resp.NS = append(resp.NS, n.Host)
		}
	} else {
		log.Println("NS lookup error:", err)
	}

	cname, err := net.LookupCNAME(hostname)
	if err == nil {
		resp.CNAME = []string{cname}
	} else {
		log.Println("CNAME lookup error:", err)
	}

	soa, err := LookupSOA(hostname)
	if err == nil {
		resp.SOA = soa
	} else {
		log.Println("SOA lookup error:", err)
	}

	_, srvs, err := net.LookupSRV("", "", hostname)
	if err == nil {
		resp.SRV = srvs
	} else {
		log.Println("SRV lookup error:", err)
	}

	ptr, err := net.LookupAddr(hostname)
	if err == nil {
		resp.PTR = ptr
	} else {
		log.Println("PTR lookup error:", err)
	}

	return resp, nil
}

func LookupSOA(hostname string) (*SOAResponse, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeSOA)
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("SOA lookup failed: %s", dns.RcodeToString[r.Rcode])
	}
	if len(r.Answer) == 0 {
		return nil, fmt.Errorf("SOA record not found for %s", hostname)
	}
	soa := r.Answer[0].(*dns.SOA)
	return &SOAResponse{
		NS:      soa.Ns,
		MBox:    soa.Mbox,
		Serial:  soa.Serial,
		Refresh: soa.Refresh,
		Retry:   soa.Retry,
		Expire:  soa.Expire,
		Minttl:  soa.Minttl,
	}, nil
}
