/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-26 16:43
**/

package app

import (
	"fmt"
	"net"

	"github.com/lemonyxk/discover/message"
)

func ParseAddr(ad string) *message.Server {
	addr, err := net.ResolveTCPAddr("tcp", ad)
	if err != nil {
		panic(err)
	}

	host, _, err := net.SplitHostPort(ad)
	if err != nil {
		panic(err)
	}

	return &message.Server{
		Addr: fmt.Sprintf("%s:%d", host, addr.Port),
		Http: fmt.Sprintf("%s:%d", host, addr.Port),
		Raft: fmt.Sprintf("%s:%d", host, addr.Port+1000),
		Tcp:  fmt.Sprintf("%s:%d", host, addr.Port+2000),
	}
}

func RaftAddr2Addr(raftAddr string) *message.Server {
	var addr, err = net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		panic(err)
	}
	host, _, err := net.SplitHostPort(raftAddr)
	if err != nil {
		panic(err)
	}
	return ParseAddr(fmt.Sprintf("%s:%d", host, addr.Port-1000))
}
