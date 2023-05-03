package utils

import "net"

func GetIP(domain string) (ip string, err error) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		return "", err
	}
	return addrs[0], nil
}
