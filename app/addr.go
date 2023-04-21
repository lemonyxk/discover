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

	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/utils"
)

func ParseAddr(ad string) *message.Server {
	host, port, err := utils.SplitHostPort(ad)
	if err != nil {
		panic(err)
	}

	return &message.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Http: fmt.Sprintf("%s:%d", host, port),
		Raft: fmt.Sprintf("%s:%d", host, port+1000),
		Tcp:  fmt.Sprintf("%s:%d", host, port+2000),
	}
}

func RaftAddr2Addr(raftAddr string) *message.Server {
	host, port, err := utils.SplitHostPort(raftAddr)
	if err != nil {
		panic(err)
	}
	return ParseAddr(fmt.Sprintf("%s:%d", host, port-1000))
}
