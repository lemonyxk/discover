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
	"bytes"
	"errors"

	"github.com/lemonyxk/kitty/socket/http/client"
	"github.com/lemonyxk/utils"
)

func (dis *discover) Get(key string) (string, error) {
	var res = client.Get(dis.url("/Get/"+key, false)).Query(nil).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	var code = utils.Json.Bytes(res.Bytes()).Get("code").Int()
	var msg = utils.Json.Bytes(res.Bytes()).Get("msg").Bytes()
	if code != 200 {
		return "", errors.New(string(msg))
	}
	return string(msg), nil
}

func (dis *discover) Set(key, value string) (string, error) {
	var buf = bytes.NewBuffer([]byte(value))
	var res = client.Post(dis.url("/Set/"+key, true)).Raw(buf).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	var code = utils.Json.Bytes(res.Bytes()).Get("code").Int()
	var msg = utils.Json.Bytes(res.Bytes()).Get("msg").Bytes()
	if code != 200 {
		return "", errors.New(string(msg))
	}
	return string(msg), nil
}

func (dis *discover) Delete(key string) (string, error) {
	var res = client.Post(dis.url("/Delete/"+key, true)).Form(nil).Send()
	if res.Error() != nil {
		return "", res.Error()
	}
	var code = utils.Json.Bytes(res.Bytes()).Get("code").Int()
	var msg = utils.Json.Bytes(res.Bytes()).Get("msg").Bytes()
	if code != 200 {
		return "", errors.New(string(msg))
	}
	return string(msg), nil
}
