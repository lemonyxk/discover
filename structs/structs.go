/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 14:34
**/

package structs

type WhoIsMaster struct {
	Addr      Address
	Timestamp int64
	IsMaster  bool
}

type Address struct {
	Addr string
	Http string
	Raft string
	Tcp  string
	Udp  string
}

type ServerInfo struct {
	ServerName string
	Addr       string
}

type Register struct {
	ServerList []string
	KeyList    []string
	ServerInfo *ServerInfo
}
