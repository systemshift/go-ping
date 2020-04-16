package main

import (
	"fmt"
	"net"
	"os"
)

func grabIP(userInput string) net.IP {
	// parse IP string into type net.IP, if nill that means userInput is a domain string
	ip := net.ParseIP(userInput)
	if ip == nil {
		ip, _ := net.LookupIP(userInput)
		return ip[0]
	}
	return ip
}

func main() {
	args := os.Args
	ip := grabIP(args[1])
	fmt.Println(ip)
}
