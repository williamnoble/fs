package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func getLocalIPAddress() (host net.IP, err error) {
	hostName, err := os.Hostname()
	if err != nil {
		return []byte{}, ErrUnableToRetrieveHostname
	}

	addrs, _ := net.LookupIP(hostName)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Printf("from ip.go : ipv4 %s\n", ipv4)
			return ipv4, nil
		}
	}
	return []byte{}, ErrUnableToRetrieveLocalIP
}

func localIPAddres() (ip string) {
	localIP, err := getLocalIPAddress()
	if err != nil {
		switch err {
		case ErrUnableToRetrieveLocalIP:
			log.Fatal(ErrGeneralisedHostError, err)
		case ErrUnableToRetrieveHostname:
			log.Fatal(ErrUnableToRetrieveHostname, err)
		default:
			log.Fatal(err) // cannot recover
		}
	}
	ip = localIP.String()
	return
}
