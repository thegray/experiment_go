package scram

import (
	"crypto/sha512"
	"hash"

	"github.com/xdg/scram"
)

var ScramSHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDGSCRAMClient struct {
	client             *scram.Client
	clientConversation *scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.clientConversation = x.client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.clientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.clientConversation.Done()
}
