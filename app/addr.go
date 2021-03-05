/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 16:43
**/

package app

import (
	"fmt"
	"net"

	"discover/structs"
)

const udpAddress = "224.0.0.250:11000"

func ParseAddr(ad string) structs.Address {

	addr, err := net.ResolveTCPAddr("tcp", ad)
	if err != nil {
		panic(err)
	}

	return structs.Address{
		Addr: addr.String(),
		Http: addr.String(),
		Raft: fmt.Sprintf("%s:%d", addr.IP, addr.Port+1000),
		Tcp:  fmt.Sprintf("%s:%d", addr.IP, addr.Port+2000),
		Udp:  udpAddress,
	}
}

func RaftAddr2Addr(raftAddr string) structs.Address {
	var addr, err = net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		panic(err)
	}

	return ParseAddr(fmt.Sprintf("%s:%d", addr.IP, addr.Port-1000))
}
