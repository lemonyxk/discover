/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-12-27 20:22
**/

package app

import (
	"sync"
)

type app struct {
	lock sync.Mutex
}
