package functions

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type SSLInfo struct {
	Certificates []*CertificateInfo `json:"certificates"`
	Error        string             `json:"error,omitempty"`
}

type CertificateInfo struct {
	Subject      string   `json:"subject"`
	Issuer       string   `json:"issuer"`
	ValidFrom    string   `json:"valid_from"`
	ValidUntil   string   `json:"valid_until"`
	SerialNumber string   `json:"serial_number"`
	SignatureAlg string   `json:"signature_algorithm"`
	KeyUsage     int      `json:"key_usage"`
	IsCACert     bool     `json:"is_ca_cert"`
	DNSNames     []string `json:"dns_names"`
}

func GetSSLInfo(domain string) (*SSLInfo, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("https://%s", domain)
	resp, err := client.Get(url)
	if err != nil {
		return &SSLInfo{Error: err.Error()}, err
	}
	defer resp.Body.Close()

	state := resp.TLS
	if state == nil {
		return &SSLInfo{}, err
	}

	certs := state.PeerCertificates
	if len(certs) == 0 {
		return &SSLInfo{}, err
	}

	var certificateList []*CertificateInfo
	for _, cert := range certs {
		certificate := &CertificateInfo{
			Subject:      cert.Subject.CommonName,
			Issuer:       cert.Issuer.CommonName,
			ValidFrom:    cert.NotBefore.Format("2006-01-02 15:04:05"),
			ValidUntil:   cert.NotAfter.Format("2006-01-02 15:04:05"),
			SerialNumber: cert.SerialNumber.String(),
			SignatureAlg: cert.SignatureAlgorithm.String(),
			KeyUsage:     int(cert.KeyUsage),
			IsCACert:     cert.IsCA,
			DNSNames:     cert.DNSNames,
		}
		certificateList = append(certificateList, certificate)
	}

	return &SSLInfo{Certificates: certificateList}, nil
}
