/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 19:59
**/

package app

type Config struct {
	Addr   string `json:"addr"`
	Http   string `json:"http"`
	Tcp    string `json:"tcp"`
	Raft   string `json:"raft"`
	Secret string `json:"secret"`
	Dir    string `json:"dir"`
	Debug  bool   `json:"debug"`
}
