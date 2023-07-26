package functions

import (
	"fmt"
	"net"
	"time"
)

var knownServices = map[int]string{
	21:   "FTP",
	22:   "SSH",
	23:   "Telnet",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	110:  "POP3",
	143:  "IMAP",
	443:  "HTTPS",
	3306: "MySQL",
}

func ScanPorts(hostname string) ([]string, error) {
	openPorts := make(chan int)
	done := make(chan struct{})
	var results []string

	for port := 0; port <= 65535; port++ {
		go func(port int) {
			address := fmt.Sprintf("%s:%d", hostname, port)
			conn, err := net.DialTimeout("tcp", address, 1*time.Second) // Timeout set to 1 second
			if err == nil {
				openPorts <- port
				conn.Close()
			}
			done <- struct{}{}
		}(port)
	}

	go func() {
		for port := 0; port <= 65535; port++ {
			<-done
		}
		close(openPorts)
	}()

	timeout := time.After(5 * time.Second)

	for {
		select {
		case port, ok := <-openPorts:
			if !ok {
				return results, nil
			}
			serviceName := knownServices[port]
			if serviceName != "" {
				results = append(results, fmt.Sprintf("%s (%d)", serviceName, port))
			} else {
				results = append(results, fmt.Sprintf("Unknown service (%d)", port))
			}
		case <-timeout:
			return nil, fmt.Errorf("scanning timeout, some ports may still be unresponsive")
		}
	}
}
