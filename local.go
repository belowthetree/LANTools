package main

import "net"

func GetIntranetIp() []string {
	address, _ := net.InterfaceAddrs()
	var res []string
	for _, addr := range address{
		if ip, ok := addr.(*net.IPNet);ok && !ip.IP.IsLoopback(){
			if ip.IP.To4() != nil{
				res = append(res, ip.IP.To4().String())
			}
		}
	}
	return res
}