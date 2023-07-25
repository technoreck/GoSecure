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

	// Concurrent port scanning for all known ports (0 to 65535)
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

	// Wait for all goroutines to finish
	go func() {
		for port := 0; port <= 65535; port++ {
			<-done
		}
		close(openPorts)
	}()

	timeout := time.After(5 * time.Second) // Total timeout set to 5 seconds

	// Collect scan results with timeout
	for {
		select {
		case port, ok := <-openPorts:
			if !ok {
				// openPorts channel closed, scanning is complete
				return results, nil
			}
			serviceName := knownServices[port]
			if serviceName != "" {
				results = append(results, fmt.Sprintf("%s (%d)", serviceName, port))
			} else {
				results = append(results, fmt.Sprintf("Unknown service (%d)", port))
			}
		case <-timeout:
			// Timeout reached
			return nil, fmt.Errorf("scanning timeout, some ports may still be unresponsive")
		}
	}
}
