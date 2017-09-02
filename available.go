package domainfinder

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const whoisServer string = "com.whois-servers.net"

type Available interface {
	Available(domain string) (bool, error)
}

type available struct {
	match string
}

func NewAvailable() Available {
	return available{
		match: "no match",
	}
}

func (a available) Available(domain string) (bool, error) {
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return false, err
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(domain + "rn")); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), a.match) {
			return false, fmt.Errorf("No match for %s", domain)
		}
	}
	return true, nil
}
