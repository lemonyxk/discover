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

type Op struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}
