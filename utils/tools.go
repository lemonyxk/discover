/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-21 12:17
**/

package utils

import (
	"net"
	"strconv"
)

func SplitHostPort(addr string) (string, int, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return host, portInt, nil
}
