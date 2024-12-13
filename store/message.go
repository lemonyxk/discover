/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-21 21:29
**/

package store

import (
	"encoding/binary"
	"errors"
	json "github.com/bytedance/sonic"
)

var ErrInvalidMessage = errors.New("invalid message")

type Op byte

const (
	Delete Op = iota
	Set
	Clear
)

type KV struct {
	Key   string                `json:"key"`
	Value json.NoCopyRawMessage `json:"value"`
}

type Message struct {
	Op    Op
	Key   string
	Value []byte
}

func Build(message *Message) []byte {
	// 0: 0 set 1 delete
	// 1: key length
	// 2-5: value length
	var kl = len(message.Key)
	var vl = len(message.Value)
	var build = make([]byte, 6+kl+vl)
	build[0] = byte(message.Op)
	build[1] = byte(kl)
	binary.BigEndian.PutUint32(build[2:6], uint32(vl))
	copy(build[6:6+kl], message.Key)
	copy(build[6+kl:], message.Value)
	return build
}

func Parse(data []byte) (*Message, error) {
	if len(data) < 6 {
		return nil, ErrInvalidMessage
	}

	var message Message
	message.Op = Op(data[0])
	var kl = int(data[1])
	var vl = int(binary.BigEndian.Uint32(data[2:6]))

	if len(data) < 6+kl+vl {
		return nil, ErrInvalidMessage
	}

	message.Key = string(data[6 : 6+kl])
	message.Value = data[6+kl : 6+kl+vl]
	return &message, nil
}

func ParseMulti(data []byte) ([]*Message, error) {
	var messages []*Message
	for len(data) > 0 {
		var message, err = Parse(data)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
		data = data[6+len(message.Key)+len(message.Value):]
	}
	return messages, nil
}
