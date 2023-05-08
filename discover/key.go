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

	"github.com/lemonyxk/kitty/errors"
	"github.com/lemonyxk/kitty/socket/http/client"
)

func (dis *Client) Get(key string) ([]byte, error) {
	var res = client.Get(dis.url("/Get/"+key, false)).Query(nil).Send()
	if res.Error() != nil {
		return nil, errors.New(res.Error())
	}

	if res.Code() != 200 {
		return nil, errors.New(res.String())
	}

	return res.Bytes(), nil
}

func (dis *Client) All() ([]byte, error) {
	var res = client.Get(dis.url("/All", false)).Query(nil).Send()
	if res.Error() != nil {
		return nil, errors.New(res.Error())
	}

	if res.Code() != 200 {
		return nil, errors.New(res.String())
	}

	return res.Bytes(), nil
}

func (dis *Client) Set(key string, value []byte) (string, error) {
	var res = client.Post(dis.url("/Set/"+key, true)).Raw(bytes.NewReader(value)).Send()
	if res.Error() != nil {
		return "", res.Error()
	}

	if res.Code() != 200 {
		return "", errors.New(res.String())
	}
	return string(res.Bytes()), nil
}

func (dis *Client) Delete(key string) (string, error) {
	var res = client.Post(dis.url("/Delete/"+key, true)).Raw(nil).Send()
	if res.Error() != nil {
		return "", res.Error()
	}

	if res.Code() != 200 {
		return "", errors.New(res.String())
	}
	return string(res.Bytes()), nil
}
