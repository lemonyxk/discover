/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-26 14:34
**/

package structs

import "github.com/lemonyxk/discover/message"

type Register struct {
	ServerList []string
	KeyList    []string
	ServerInfo *message.ServerInfo
}
