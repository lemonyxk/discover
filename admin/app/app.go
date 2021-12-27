/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-12-27 20:22
**/

package app

import (
	"sync"

	"github.com/lemoyxk/discover/message"
)

var App = &app{
	ServerMap: &serverMap{
		servers: make(map[string]*message.WhoIsMaster),
	},
}

type app struct {
	ServerMap *serverMap
	lock      sync.Mutex
}
