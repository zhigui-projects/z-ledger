/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"bytes"
	proto "github.com/hyperledger/fabric-protos-go/gossip"
	"github.com/hyperledger/fabric/gossip/common"
	"github.com/hyperledger/fabric/gossip/protoext"
	"sync"
)

type gossip interface {
	// PeersOfChannel returns the NetworkMembers considered alive in a channel
	//PeersOfChannel(channel common.ChannelID) []discovery.NetworkMember

	// Accept returns a dedicated read-only channel for messages sent by other nodes that match a certain predicate.
	// If passThrough is false, the messages are processed by the gossip layer beforehand.
	// If passThrough is true, the gossip layer doesn't intervene and the messages
	// can be used to send a reply back to the sender
	Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan protoext.ReceivedMessage)

	// Gossip sends a message to other peers to the network
	Gossip(msg *proto.GossipMessage)

	Stop()
}

type gossipImpl struct {
	gossip   gossip
	channel  common.ChannelID
	doneCh   chan struct{}
	stopOnce *sync.Once
}

func NewGossipImpl(g gossip, channel string) ArchiveGossip {
	return &gossipImpl{
		gossip:   g,
		channel:  common.ChannelID(channel),
		doneCh:   make(chan struct{}),
		stopOnce: &sync.Once{},
	}
}

func (g *gossipImpl) Accept() <-chan *proto.GossipMessage {
	archiveCh, _ := g.gossip.Accept(func(message interface{}) bool {
		// Get only archive messages
		return message.(*proto.GossipMessage).Tag == proto.GossipMessage_CHAN_ONLY &&
			protoext.IsArchiveMsg(message.(*proto.GossipMessage)) &&
			bytes.Equal(message.(*proto.GossipMessage).Channel, g.channel)
	}, false)

	msgCh := make(chan *proto.GossipMessage)

	go func(inCh <-chan *proto.GossipMessage, outCh chan *proto.GossipMessage, stopCh chan struct{}) {
		for {
			select {
			case <-stopCh:
				return
			case gossipMsg, ok := <-inCh:
				if ok {
					outCh <- gossipMsg
				} else {
					return
				}
			}
		}
	}(archiveCh, msgCh, g.doneCh)
	return msgCh
}

func (g *gossipImpl) Gossip(msg *proto.GossipMessage) {
	g.gossip.Gossip(msg)
}

func (g *gossipImpl) Stop() {
	g.stopOnce.Do(func() {
		close(g.doneCh)
	})
}
