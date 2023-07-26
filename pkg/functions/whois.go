package functions

import (
	"bufio"
	"net"
	"strings"
)

const whoisServer = "whois.iana.org"

// WHOISInfo holds the WHOIS information for a domain.
type WHOISInfo struct {
	Data string `json:"data"`
}

// GetWHOISInfo retrieves the WHOIS information for the given domain.
func GetWHOISInfo(domain string) (*WHOISInfo, error) {
	info, err := getWHOISInfo(domain)
	if err != nil {
		return nil, err
	}

	// Remove the unwanted sections from the WHOIS information
	info = removeUnwantedSections(info)

	return &WHOISInfo{Data: info}, nil
}

func getWHOISInfo(domain string) (string, error) {
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.Write([]byte(domain + "\r\n"))

	var result strings.Builder
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		result.WriteString(line + "\n")

		if strings.HasPrefix(line, "refer:") {
			refServer := strings.TrimSpace(strings.TrimPrefix(line, "refer:"))
			if refServer != "" {
				return getWHOISInfoFromRefServer(domain, refServer)
			}
		}
	}

	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return result.String(), nil
}

func getWHOISInfoFromRefServer(domain, refServer string) (string, error) {
	conn, err := net.Dial("tcp", refServer+":43")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.Write([]byte(domain + "\r\n"))

	var result strings.Builder
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		result.WriteString(line + "\n")
	}

	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return result.String(), nil
}

func removeUnwantedSections(info string) string {
	// Define the sections to be removed
	unwantedSections := []string{
		"URL of the ICANN Whois Inaccuracy Complaint Form:",
		">>> Last update of whois database:",
		"For more information on Whois status codes, please visit",
		"NOTICE: The expiration date displayed in this record is the date the",
		"database through the use of electronic processes that are high-volume and",
		"automated except as reasonably necessary to register domain names or",
		"modify existing registrations; the Data in VeriSign Global Registry",
		"Services' (\"VeriSign\") Whois database is provided by VeriSign for",
		"information purposes only, and to assist persons in obtaining information",
		"about or related to a domain name registration record. VeriSign does not",
		"guarantee its accuracy. By submitting a Whois query, you agree to abide",
		"by the following terms of use: You agree that you may use this Data only",
		"for lawful purposes and that under no circumstances will you use this Data",
		"to: (1) allow, enable, or otherwise support the transmission of mass",
		"unsolicited, commercial advertising or solicitations via e-mail, telephone,",
		"or facsimile; or (2) enable high volume, automated, electronic processes",
		"that apply to VeriSign (or its computer systems). The compilation,",
		"repackaging, dissemination or other use of this Data is expressly",
		"prohibited without the prior written consent of VeriSign. You agree not to",
		"use electronic processes that are automated and high-volume to access or",
		"query the Whois database except as reasonably necessary to register",
		"domain names or modify existing registrations. VeriSign reserves the right",
		"to restrict your access to the Whois database in its sole discretion to ensure",
		"operational stability.  VeriSign may restrict or terminate your access to the",
		"Whois database for failure to abide by these terms of use. VeriSign",
		"reserves the right to modify these terms at any time.",
		"The Registry database contains ONLY .COM, .NET, .EDU domains and",
		"Registrars.",
	}

	lines := strings.Split(info, "\n")
	var filteredLines []string

	// Iterate through the lines and exclude the unwanted sections
	for _, line := range lines {
		includeLine := true
		for _, unwantedSection := range unwantedSections {
			if strings.Contains(line, unwantedSection) {
				includeLine = false
				break
			}
		}

		if includeLine {
			filteredLines = append(filteredLines, line)
		}
	}

	return strings.Join(filteredLines, "\n")
}
