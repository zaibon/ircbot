package ircbot

import (
	"encoding/json"
)

type IrcRequest struct {
	Command string
	Channel string
	Args    []string
}

func NewIrcRequest() *IrcRequest {
	return &IrcRequest{
		Args: make([]string, 1),
	}
}

func EncodeIrcReq(req *IrcRequest) ([]byte, error) {
	bob, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return bob, nil
}

func DecodeIrcReq(bob []byte) (*IrcRequest, error) {
	req := NewIrcRequest()

	err := json.Unmarshal(bob, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
