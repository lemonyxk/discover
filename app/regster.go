/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-03-04 10:35
**/

package app

import (
	"sync"

	"discover/structs"
)

type register struct {
	mux  sync.Mutex
	data map[int64]*structs.Register
}

func (r *register) All() map[int64]*structs.Register {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.data
}

func (r *register) Get(fd int64) *structs.Register {
	r.mux.Lock()
	defer r.mux.Unlock()
	var list, ok = r.data[fd]
	if !ok {
		return nil
	}
	return list
}

func (r *register) Set(fd int64, register *structs.Register) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.data[fd] = register
	return true
}

func (r *register) Delete(fd int64) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	var _, ok = r.data[fd]
	if !ok {
		return false
	}

	delete(r.data, fd)

	return true
}

func (r *register) Destroy() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.data = map[int64]*structs.Register{}
}
