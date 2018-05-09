package dht

import (
	"context"
	"sync"
	"time"

	logrus "github.com/sirupsen/logrus"
)

const numPeersNear int = 15

type QueryType int

const (
	PeerInfoQuery QueryType = iota
	ProviderQuery
	ValueQuery
)

type query struct {
	dht              *DHT
	id               string
	key              string
	queryType        QueryType
	closestPeerID    string
	contactedPeers   sync.Map
	incomingMessages chan interface{}
	outgoingMessages chan interface{}
	// lock             *sync.RWMutex
}

func (q *query) Run(ctx context.Context) {
	logger := logrus.WithField("resp", q.key)

	go func() {
		// send what we know about the key
		switch q.queryType {
		case PeerInfoQuery:
			if peerInfo, err := q.dht.registry.GetPeerInfo(q.key); err != nil {
				q.outgoingMessages <- peerInfo
			}
		case ProviderQuery:
			if providers, err := q.dht.store.GetProviders(q.key); err != nil {
				for _, provider := range providers {
					q.outgoingMessages <- provider
				}
			}
		case ValueQuery:
			value, err := q.dht.store.GetValue(q.key)
			if err != nil {
				break
			}
			q.outgoingMessages <- value
		}

		// and now, wait for something to happen
		for {
			select {
			case incomingMessage := <-q.incomingMessages:
				logger.Debug("Processing incoming message")
				switch message := incomingMessage.(type) {
				case *messagePutPeerInfo:
					q.outgoingMessages <- &message.PeerInfo
					q.nextIfCloser(message.PeerID)
				case *messagePutProviders:
					q.outgoingMessages <- message.PeerIDs
					q.nextIfCloser(message.SenderPeerInfo.ID)
				case *messagePutValue:
					q.outgoingMessages <- message.Value
					q.nextIfCloser(message.SenderPeerInfo.ID)
				}

			case <-time.After(maxQueryTime):
				close(q.outgoingMessages)
				return

			case <-ctx.Done():
				close(q.outgoingMessages)
				return
			}
		}
	}()

	// start looking for the node
	go q.next()
}

func (q *query) nextIfCloser(newPeerID string) {
	if q.closestPeerID == "" {
		q.closestPeerID = newPeerID
		q.next()
	} else {
		if comparePeers(q.closestPeerID, newPeerID, q.key) == newPeerID {
			q.closestPeerID = newPeerID
			q.next()
		}
	}
}

func (q *query) next() {
	// find closest peers
	closestPeers, err := q.dht.FindPeersClosestTo(q.key, numPeersNear)
	if err != nil {
		logrus.WithError(err).Error("Failed find peers near")
		return
	}

	peersToAsk := []string{}
	for _, peerInfo := range closestPeers {
		// skip the ones we've already asked
		if _, ok := q.contactedPeers.Load(peerInfo.ID); ok {
			continue
		}
		peersToAsk = append(peersToAsk, peerInfo.ID)
		q.contactedPeers.Store(peerInfo.ID, true)
	}

	var payloadType string
	var req interface{}

	switch q.queryType {
	case PeerInfoQuery:
		payloadType = PayloadTypeGetPeerInfo
		req = messageGetPeerInfo{
			SenderPeerInfo: *q.dht.registry.GetLocalPeerInfo(),
			RequestID:      q.id,
			PeerID:         q.key,
		}
	case ProviderQuery:
		payloadType = PayloadTypeGetProviders
		req = messageGetProviders{
			SenderPeerInfo: *q.dht.registry.GetLocalPeerInfo(),
			RequestID:      q.id,
			Key:            q.key,
		}
	case ValueQuery:
		payloadType = PayloadTypeGetValue
		req = messageGetValue{
			SenderPeerInfo: *q.dht.registry.GetLocalPeerInfo(),
			RequestID:      q.id,
			Key:            q.key,
		}
	default:
		return
	}

	ctx := context.Background()
	q.dht.wire.Send(ctx, wireExtention, payloadType, req, peersToAsk)
}