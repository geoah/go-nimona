package net

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	igd "github.com/emersion/go-upnp-igd"
	ucodec "github.com/ugorji/go/codec"
	"go.uber.org/zap"

	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
	"nimona.io/go/log"
	"nimona.io/go/peers"
)

// Networker interface for mocking Network
type Networker interface {
	Dial(ctx context.Context, address string) (*Connection, error)
	Listen(ctx context.Context, addrress string) (chan *Connection, error)
}

// NewNetwork creates a new p2p network using an address book
func NewNetwork(addressBook *peers.AddressBook) (*Network, error) {
	return &Network{
		addressBook: addressBook,
	}, nil
}

// Network allows dialing and listening for p2p connections
type Network struct {
	addressBook *peers.AddressBook
}

// Dial to a peer and return a net.Conn or error
func (n *Network) Dial(ctx context.Context, address string) (*Connection, error) {
	logger := log.Logger(ctx)

	addressType := strings.Split(address, ":")[0]
	switch addressType {
	case "peer":
		logger.Debug("dialing peer", zap.String("peer", address))
		peerID := strings.Replace(address, "peer:", "", 1)
		peerInfo, err := n.addressBook.GetPeerInfo(peerID)
		if err != nil {
			return nil, err
		}

		if len(peerInfo.Addresses) == 0 {
			return nil, ErrNoAddresses
		}

		for _, addr := range peerInfo.Addresses {
			conn, err := n.Dial(ctx, addr)
			if err == nil {
				return conn, nil
			}
		}

		return nil, ErrAllAddressesFailed

	case "tcp":
		addr := strings.Replace(address, "tcp:", "", 1)
		dialer := net.Dialer{Timeout: time.Second}
		logger.Debug("dialing", zap.String("address", addr))
		tcpConn, err := dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, err
		}

		if tcpConn == nil {
			return nil, ErrAllAddressesFailed
		}

		conn := &Connection{
			Conn:     tcpConn,
			RemoteID: "", // we don't really know who the other side is
		}

		signer := n.addressBook.GetLocalPeerKey()
		nonce := RandStringBytesMaskImprSrc(8)
		syn := &HandshakeSyn{
			Nonce:    nonce,
			PeerInfo: n.addressBook.GetLocalPeerInfo(),
		}
		sig, err := crypto.Sign(syn, signer)
		if err != nil {
			return nil, err
		}

		syn.Signature = sig

		if err := Write(syn, conn); err != nil {
			return nil, err
		}

		synAckIf, err := Read(conn)
		if err != nil {
			return nil, err
		}

		synAck, ok := synAckIf.(*HandshakeSynAck)
		if !ok {
			return nil, errors.New("unexpected syn ack block")
		}

		if synAck.Nonce != nonce {
			return nil, errors.New("invalid handhshake.syn-ack")
		}

		// store who is on the other side - peer id
		conn.RemoteID = synAck.Signature.Key.Thumbprint()
		if err := n.addressBook.PutPeerInfo(synAck.PeerInfo); err != nil {
			log.DefaultLogger.Panic("could not add remote peer", zap.Error(err))
		}

		ack := &HandshakeAck{
			Nonce: nonce,
		}
		sig, err = crypto.Sign(ack, signer)
		if err != nil {
			return nil, err
		}

		ack.Signature = sig

		if err := Write(ack, conn); err != nil {
			return nil, err
		}

		return conn, nil
	default:
		logger.Info("not sure how to dial", zap.String("address", address), zap.String("type", addressType))
	}

	return nil, ErrNoAddresses
}

// Listen on an address
// TODO do we need to return a listener?
func (n *Network) Listen(ctx context.Context, address string) (chan *Connection, error) {
	logger := log.Logger(ctx).Named("network")
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	port := tcpListener.Addr().(*net.TCPAddr).Port
	logger.Info("Listening and service nimona", zap.Int("port", port))
	addresses := GetAddresses(tcpListener)
	devices := make(chan igd.Device, 10)

	if n.addressBook.LocalHostname != "" {
		addresses = append(addresses, fmtAddress(n.addressBook.LocalHostname, port))
	}

	upnp := true
	upnpFlag := os.Getenv("UPNP")
	if upnpFlag != "" {
		upnp, _ = strconv.ParseBool(upnpFlag)
	}
	if upnp {
		logger.Info("Trying to find external IP and open port")
		go func() {
			if err := igd.Discover(devices, 2*time.Second); err != nil {
				logger.Error("could not discover devices", zap.Error(err))
			}
		}()
		for device := range devices {
			externalAddress, err := device.GetExternalIPAddress()
			if err != nil {
				logger.Error("could not get external ip", zap.Error(err))
				continue
			}
			desc := "nimona"
			ttl := time.Hour * 24 * 365
			if _, err := device.AddPortMapping(igd.TCP, port, port, desc, ttl); err != nil {
				logger.Error("could not add port mapping", zap.Error(err))
			} else {
				addresses = append(addresses, fmtAddress(externalAddress.String(), port))
			}
		}
	}

	logger.Info("Started listening", zap.Strings("addresses", addresses))
	n.addressBook.AddLocalPeerAddress(addresses...)

	cconn := make(chan *Connection, 10)
	go func() {
		signer := n.addressBook.GetLocalPeerKey()
		for {
			tcpConn, err := tcpListener.Accept()
			if err != nil {
				log.DefaultLogger.Warn("could not accept connection", zap.Error(err))
				// TODO close conn?
				return
			}

			conn := &Connection{
				Conn:     tcpConn,
				RemoteID: "unknown: handshaking",
			}

			synIf, err := Read(conn)
			if err != nil {
				log.DefaultLogger.Warn("waiting for syn failed", zap.Error(err))
				// TODO close conn?
				continue
			}

			syn, ok := synIf.(*HandshakeSyn)
			if !ok {
				// TODO close conn?
				continue
			}

			// TODO check type

			nonce := syn.Nonce

			// store the peer on the other side
			if err := n.addressBook.PutPeerInfo(syn.PeerInfo); err != nil {
				log.DefaultLogger.Panic("could not add remote peer", zap.Error(err))
			}

			synAck := &HandshakeSynAck{
				Nonce:    nonce,
				PeerInfo: n.addressBook.GetLocalPeerInfo(),
			}
			sig, err := crypto.Sign(synAck, signer)
			if err != nil {
				log.DefaultLogger.Warn("could not sigh for syn ack block", zap.Error(err))
				// TODO close conn?
				continue
			}
			synAck.Signature = sig
			if err := Write(synAck, conn); err != nil {
				log.DefaultLogger.Warn("sending for syn-ack failed", zap.Error(err))
				// TODO close conn?
				continue
			}

			ackIf, err := Read(conn)
			if err != nil {
				log.DefaultLogger.Warn("waiting for ack failed", zap.Error(err))
				// TODO close conn?
				continue
			}

			ack, ok := ackIf.(*HandshakeAck)
			if !ok {
				// TODO close conn?
				continue
			}

			if ack.Nonce != nonce {
				log.DefaultLogger.Warn("validating syn to ack nonce failed")
				// TODO close conn?
				continue
			}

			conn.RemoteID = ack.Signature.Key.Thumbprint()
			cconn <- conn
		}
	}()

	return cconn, nil
}

func Write(p interface{}, conn *Connection) error {
	conn.Conn.SetWriteDeadline(time.Now().Add(time.Second))
	if p == nil {
		log.DefaultLogger.Error("block for fw cannot be nil")
		return errors.New("missing block")
	}

	b, err := encoding.Marshal(p)
	if err != nil {
		return err
	}

	if _, err := conn.Conn.Write(b); err != nil {
		return err
	}

	SendBlockEvent(
		"incoming",
		encoding.GetType(p),
		len(b),
	)

	if os.Getenv("DEBUG_BLOCKS") == "true" {
		m := map[string]interface{}{}
		encoding.UnmarshalSimple(b, &m)
		b, _ := json.MarshalIndent(m, "", "  ")
		log.DefaultLogger.Info(string(b), zap.String("remoteID", conn.RemoteID), zap.String("direction", "outgoing"))
	}
	return nil
}

func Read(conn *Connection) (interface{}, error) {
	logger := log.DefaultLogger

	pDecoder := ucodec.NewDecoder(conn.Conn, encoding.RawCborHandler())
	r := ucodec.Raw{}
	if err := pDecoder.Decode(&r); err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			spew.Dump(r)
			logger.Error("Recovered while processing", zap.Any("r", r))
		}
	}()

	v, err := encoding.Unmarshal(r)
	if err != nil {
		return nil, err
	}

	// TODO(geoah) fix sig
	// if b.Signature != nil {
	// 	if err := crypto.Verify(b.Signature, d); err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	fmt.Println("--------------------------------------------------------")
	// 	fmt.Println("----- BLOCK NOT SIGNED ---------------------------------")
	// 	fmt.Println("--------------------------------------------------------")
	// 	fmt.Println("-----", b.Type)
	// 	fmt.Println("-----", b.Payload)
	// 	fmt.Println("--------------------------------------------------------")
	// }

	SendBlockEvent(
		"incoming",
		encoding.GetType(v),
		pDecoder.NumBytesRead(),
	)
	if os.Getenv("DEBUG_BLOCKS") == "true" {
		b, _ := encoding.Marshal(v)
		m := map[string]interface{}{}
		encoding.UnmarshalSimple(b, &m)
		b, _ = json.MarshalIndent(m, "", "  ")
		logger.Info(string(b), zap.String("remoteID", conn.RemoteID), zap.String("direction", "incoming"))
	}
	return v, nil
}
