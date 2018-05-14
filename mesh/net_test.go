package mesh

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNet(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "nimona")
	assert.NoError(t, err)

	tempFile1 := path.Join(tempDir, "key1")
	key1, err := LoadOrCreatePrivateKey(tempFile1)
	localPeerID := Thumbprint(key1)
	assert.NoError(t, err)

	tempFile2 := path.Join(tempDir, "key2")
	key2, err := LoadOrCreatePrivateKey(tempFile2)
	remotePeerID := Thumbprint(key2)
	assert.NoError(t, err)

	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)

	localHandled := false
	remoteHandled := false
	handler := func(conn net.Conn) net.Conn {
		peerID := conn.LocalAddr().String()
		fmt.Println("hit handler, local address", peerID)
		if peerID == localPeerID {
			localHandled = true
		} else {
			remoteHandled = true
		}
		wg.Done()
		return conn
	}

	localRegistry := NewRegisty(key1)
	localHandler := &MockHandler{}
	localHandler.On("Initiate", mock.Anything).Return(handler, nil)
	localNet := New(localRegistry)
	localNet.handlers["hi"] = localHandler

	remoteRegistry := NewRegisty(key2)
	remoteHandler := &MockHandler{}
	remoteHandler.On("Handle", mock.Anything).Return(handler, nil)
	remoteNet := New(remoteRegistry)
	remoteNet.handlers["hi"] = remoteHandler

	_, _, localListenErr := localNet.Listen("127.0.0.1:0")
	_, remoteAddr, remoteListenErr := remoteNet.Listen("127.0.0.1:0")
	assert.NoError(t, localListenErr)
	assert.NoError(t, remoteListenErr)

	localRegistry.PutPeerInfo(&PeerInfo{
		ID:        remotePeerID,
		Addresses: []string{remoteAddr},
	})

	conn, err := localNet.Dial(ctx, remotePeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()

	assert.Equal(t, localPeerID, conn.LocalAddr().String())
	assert.Equal(t, remotePeerID, conn.RemoteAddr().String())

	assert.True(t, localHandled)
	assert.True(t, remoteHandled)

	localNet.Close()
	remoteNet.Close()
}

func TestReusableNet(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "nimona")
	assert.NoError(t, err)

	tempFile1 := path.Join(tempDir, "key1")
	key1, err := LoadOrCreatePrivateKey(tempFile1)
	localPeerID := Thumbprint(key1)
	assert.NoError(t, err)

	tempFile2 := path.Join(tempDir, "key2")
	key2, err := LoadOrCreatePrivateKey(tempFile2)
	remotePeerID := Thumbprint(key2)
	assert.NoError(t, err)

	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)

	localHandled := false
	remoteHandled := false
	handler := func(conn net.Conn) net.Conn {
		peerID := conn.LocalAddr().String()
		fmt.Println("hit handler, local address", peerID)
		if peerID == localPeerID {
			localHandled = true
		} else {
			remoteHandled = true
		}
		wg.Done()
		return conn
	}

	localRegistry := NewRegisty(key1)
	localHandler := &MockHandler{}
	localHandler.On("Initiate", mock.Anything).Return(handler, nil)
	localNet := New(localRegistry)
	localNet.handlers["hi"] = localHandler

	remoteRegistry := NewRegisty(key2)
	remoteHandler := &MockHandler{}
	remoteHandler.On("Handle", mock.Anything).Return(handler, nil)
	remoteNet := New(remoteRegistry)
	remoteNet.handlers["hi"] = remoteHandler

	_, _, localListenErr := localNet.Listen("127.0.0.1:0")
	_, remoteAddr, remoteListenErr := remoteNet.Listen("127.0.0.1:0")
	assert.NoError(t, localListenErr)
	assert.NoError(t, remoteListenErr)

	localRegistry.PutPeerInfo(&PeerInfo{
		ID:        remotePeerID,
		Addresses: []string{remoteAddr},
	})

	conn, err := localNet.Dial(ctx, remotePeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()

	assert.Equal(t, localPeerID, conn.LocalAddr().String())
	assert.Equal(t, remotePeerID, conn.RemoteAddr().String())

	assert.True(t, localHandled)
	assert.True(t, remoteHandled)

	localNet.Close()
	remoteNet.Close()
}

func TestReusableRedialNet(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "nimona")
	assert.NoError(t, err)

	tempFile1 := path.Join(tempDir, "key1")
	key1, err := LoadOrCreatePrivateKey(tempFile1)
	localPeerID := Thumbprint(key1)
	assert.NoError(t, err)

	tempFile2 := path.Join(tempDir, "key2")
	key2, err := LoadOrCreatePrivateKey(tempFile2)
	remotePeerID := Thumbprint(key2)
	assert.NoError(t, err)

	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)

	var localHandled int32
	var remoteHandled int32

	handler := func(conn net.Conn) net.Conn {
		peerID := conn.LocalAddr().String()
		fmt.Println("hit handler, local address", peerID)
		if peerID == localPeerID {
			fmt.Println("> HI")
			atomic.AddInt32(&localHandled, 1)
		} else {
			fmt.Println("< HI")
			atomic.AddInt32(&remoteHandled, 1)
		}
		wg.Done()
		return conn
	}

	localRegistry := NewRegisty(key1)
	localHandler := &MockHandler{}
	localHandler.On("Initiate", mock.Anything).Return(handler, nil)
	localNet := New(localRegistry)
	localNet.handlers["hi"] = localHandler

	remoteRegistry := NewRegisty(key2)
	remoteHandler := &MockHandler{}
	remoteHandler.On("Handle", mock.Anything).Return(handler, nil)
	remoteNet := New(remoteRegistry)
	remoteNet.handlers["hi"] = remoteHandler

	_, _, localListenErr := localNet.Listen("127.0.0.1:0")
	_, remoteAddr, remoteListenErr := remoteNet.Listen("127.0.0.1:0")
	assert.NoError(t, localListenErr)
	assert.NoError(t, remoteListenErr)

	localRegistry.PutPeerInfo(&PeerInfo{
		ID:        remotePeerID,
		Addresses: []string{remoteAddr},
	})

	conn, err := localNet.Dial(ctx, remotePeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()
	wg.Add(2)

	assert.Equal(t, localPeerID, conn.LocalAddr().String())
	assert.Equal(t, remotePeerID, conn.RemoteAddr().String())

	assert.Equal(t, 1, int(localHandled))
	assert.Equal(t, 1, int(remoteHandled))

	conn, err = localNet.Dial(ctx, remotePeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()

	assert.Equal(t, 2, int(localHandled))
	assert.Equal(t, 2, int(remoteHandled))

	localNet.Close()
	remoteNet.Close()

}

func TestReusableRedialRemoteNet(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "nimona")
	assert.NoError(t, err)

	tempFile1 := path.Join(tempDir, "key1")
	key1, err := LoadOrCreatePrivateKey(tempFile1)
	localPeerID := Thumbprint(key1)
	assert.NoError(t, err)

	tempFile2 := path.Join(tempDir, "key2")
	key2, err := LoadOrCreatePrivateKey(tempFile2)
	remotePeerID := Thumbprint(key2)
	assert.NoError(t, err)

	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)

	var localHandled int32
	var remoteHandled int32

	handler := func(conn net.Conn) net.Conn {
		peerID := conn.LocalAddr().String()
		fmt.Println("hit handler, local address", peerID)
		if peerID == localPeerID {
			fmt.Println("> HI")
			atomic.AddInt32(&localHandled, 1)
		} else {
			fmt.Println("< HI")
			atomic.AddInt32(&remoteHandled, 1)
		}
		wg.Done()
		return conn
	}

	localRegistry := NewRegisty(key1)
	localHandler := &MockHandler{}
	localHandler.On("Initiate", mock.Anything).Return(handler, nil)
	localHandler.On("Handle", mock.Anything).Return(handler, nil)
	localNet := New(localRegistry)
	localNet.handlers["hi"] = localHandler

	remoteRegistry := NewRegisty(key2)
	remoteHandler := &MockHandler{}
	remoteHandler.On("Handle", mock.Anything).Return(handler, nil)
	remoteHandler.On("Initiate", mock.Anything).Return(handler, nil)
	remoteNet := New(remoteRegistry)
	remoteNet.handlers["hi"] = remoteHandler

	_, localAddr, localListenErr := localNet.Listen("127.0.0.1:0")
	_, remoteAddr, remoteListenErr := remoteNet.Listen("127.0.0.1:0")
	assert.NoError(t, localListenErr)
	assert.NoError(t, remoteListenErr)

	remoteRegistry.PutPeerInfo(&PeerInfo{
		ID:        localPeerID,
		Addresses: []string{localAddr},
	})

	localRegistry.PutPeerInfo(&PeerInfo{
		ID:        remotePeerID,
		Addresses: []string{remoteAddr},
	})

	conn, err := localNet.Dial(ctx, remotePeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()
	wg.Add(2)

	assert.Equal(t, localPeerID, conn.LocalAddr().String())
	assert.Equal(t, remotePeerID, conn.RemoteAddr().String())

	assert.Equal(t, 1, int(localHandled))
	assert.Equal(t, 1, int(remoteHandled))

	conn, err = remoteNet.Dial(ctx, localPeerID, "hi")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	wg.Wait()

	assert.Equal(t, 2, int(localHandled))
	assert.Equal(t, 2, int(remoteHandled))

	localNet.Close()
	remoteNet.Close()
}
