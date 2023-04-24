/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-19 15:14
**/

package message

type ServerInfo struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type AliveResponse struct {
	Name           string        `json:"name"`
	ServerInfoList []*ServerInfo `json:"server_info_list"`
}

type Address struct {
	*Server `json:",inline"`
	Master  bool `json:"master"`
}

type Server struct {
	Addr string `json:"addr"`
	Http string `json:"http"`
	Raft string `json:"raft"`
	Tcp  string `json:"tcp"`
}
