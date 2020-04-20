package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func grabIP(userInput string) net.IPAddr {
	// parse IP or domain string into type net.IPAddr
	ip, err := net.ResolveIPAddr("ip4", userInput)
	if err != nil {
		panic(err)
		return *ip
	}
	return *ip
}

func echo(address string) (net.IPAddr, time.Duration, error) {
	ip := grabIP(address)
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return ip, 0, err
	}
	defer conn.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getegid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}

	body, err := message.Marshal(nil)
	if err != nil {
		return ip, 0, err
	}

	start_time := time.Now()

	n, err := conn.WriteTo(body, &ip)
	if err != nil {
		return ip, 0, err
	}

	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return ip, 0, err
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return ip, 0, err
	}
	end_time := time.Since(start_time)

	rm, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return ip, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return ip, end_time, nil
	default:
		return ip, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)

	}
}

func main() {
	args := os.Args
	ip := grabIP(args[1])
	fmt.Println(ip)

	dst, dur, err := echo(args[1])

	log.Printf("ping %s (%s): %s", dst, dur, err)

}
