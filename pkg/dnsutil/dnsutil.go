package dnsutil

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

//Ping send a ping to the DNS so it can record your service and IP
func Ping(address string, serviceName string) net.Addr {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
	}
	defer conn.Close()
	fmt.Fprintf(conn, "recordAddress="+serviceName)
	ip := conn.LocalAddr()
	return ip
}

//GetMyIP returns the caller's ip address by sending a blank request to google's DNS server
//and retrieving the local address from the response
func GetMyIP() net.Addr {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("udp", "8.8.8.8:80")
		if err == nil {
			break
		}
	}
	defer conn.Close()
	ip := conn.LocalAddr()
	return ip
}

//GetServiceIP queries the dns server for a currently running service and returns the IP address
func GetServiceIP(dnsAddr string, serviceName string) string {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", dnsAddr)
		if err == nil {
			break
		}
	}
	defer conn.Close()
	fmt.Fprint(conn, "getAddress="+serviceName)

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	//trim te null characters from the buffer and convert to string
	bufferText := string(bytes.Trim(buffer, "\x00"))
	return bufferText
}

//TrimPort converts an IP address into a string containing the IP without the port attached
func TrimPort(ip net.Addr) string {
	stringIP := ip.String()
	noPort := strings.Split(stringIP, ":")[0]
	return noPort
}
