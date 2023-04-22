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

	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/kitty/errors"
	"github.com/lemonyxk/kitty/socket/http/client"
)

func (dis *Client) Get(key string) (*store.Message, error) {
	var res = client.Get(dis.url("/Get/"+key, false)).Query(nil).Send()
	if res.Error() != nil {
		return nil, errors.New(res.Error())
	}

	if res.Code() != 200 {
		return nil, errors.New(res.String())
	}

	parse, err := store.Parse(res.Bytes())
	if err != nil {
		return nil, err
	}

	return parse, nil
}

func (dis *Client) All() ([]*store.Message, error) {
	var res = client.Get(dis.url("/All", false)).Query(nil).Send()
	if res.Error() != nil {
		return nil, errors.New(res.Error())
	}

	if res.Code() != 200 {
		return nil, errors.New(res.String())
	}

	parse, err := store.ParseMulti(res.Bytes())
	if err != nil {
		return nil, err
	}

	return parse, nil
}

func (dis *Client) Set(key string, value []byte) (string, error) {
	var message = &store.Message{
		Op:    store.Set,
		Key:   key,
		Value: value,
	}

	var build = store.Build(message)
	var buf = bytes.NewBuffer(build)
	var res = client.Post(dis.url("/Set/"+message.Key, true)).Raw(buf).Send()
	if res.Error() != nil {
		return "", res.Error()
	}

	if res.Code() != 200 {
		return "", errors.New(res.String())
	}
	return string(res.Bytes()), nil
}

func (dis *Client) Delete(key string) (string, error) {
	var message = &store.Message{
		Op:    store.Delete,
		Key:   key,
		Value: nil,
	}

	var build = store.Build(message)
	var buf = bytes.NewBuffer(build)
	var res = client.Post(dis.url("/Delete/"+key, true)).Raw(buf).Send()
	if res.Error() != nil {
		return "", res.Error()
	}

	if res.Code() != 200 {
		return "", errors.New(res.String())
	}
	return string(res.Bytes()), nil
}
