/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-03-05 15:40
**/

package discover

import (
	"errors"
	"strings"

	"github.com/lemonyxk/kitty/kitty"
	"github.com/lemonyxk/kitty/socket/http/client"
)

func (dis *discover) Get(key string) (string, error) {
	var res = client.Get(dis.url("/Get", false)).Query(kitty.M{"key": key}).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	if !strings.HasPrefix(res.String(), "OK") {
		return "", errors.New(res.String())
	}
	return res.String(), nil
}

func (dis *discover) Set(key, value string) (string, error) {
	var res = client.Post(dis.url("/Set", true)).Form(kitty.M{"key": key, "value": value}).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	if !strings.HasPrefix(res.String(), "OK") {
		return "", errors.New(res.String())
	}
	return res.String(), nil
}

func (dis *discover) Delete(key string) (string, error) {
	var res = client.Post(dis.url("/Delete", true)).Form(kitty.M{"key": key}).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	if !strings.HasPrefix(res.String(), "OK") {
		return "", errors.New(res.String())
	}
	return res.String(), nil
}
