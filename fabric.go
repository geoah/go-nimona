package fabric

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrNoTransport      = errors.New("Could not dial with available transports")
	ErrNoSuchMiddleware = errors.New("No such middleware")

	ContextKeyAddressPart = contextKey("addrpart")
)

func New() *Fabric {
	return &Fabric{
		transports:  map[string]Transport{},
		negotiators: map[string]NegotiatorFunc{},
		routes:      map[string]HandlerFunc{},
	}
}

type Fabric struct {
	// handlers    []Handler
	negotiators map[string]NegotiatorFunc
	transports  map[string]Transport
	routes      map[string]HandlerFunc
}

func (f *Fabric) AddTransport(n string, tr Transport) error {
	f.transports[n] = tr
	return nil
}

func (f *Fabric) AddHandlerFunc(r string, hf HandlerFunc) error {
	f.routes[r] = hf
	return nil
}

func (f *Fabric) AddNegotiatorFunc(n string, ng NegotiatorFunc) error {
	f.negotiators[n] = ng
	return nil
}

func (f *Fabric) DialContext(ctx context.Context, addr string) (Conn, error) {
	// TODO validate the address

	// figure out if the addr can be dialed and connect to the target
	c, err := f.dialTransport(ctx, addr)
	if err != nil {
		return nil, err
	}

	// handshake
	if err := f.handshake(c); err != nil {
		return nil, err
	}

	// go throught all the protocols that are defined in the address
	if err := f.Next(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (f *Fabric) dialTransport(ctx context.Context, addr string) (*conn, error) {
	st := strings.Split(addr, "/")
	ns := st[0]
	np := strings.Split(ns, ":")
	pr := np[0]

	// get protocol
	tr, err := f.getTransport(pr)
	if err != nil {
		return nil, ErrNoTransport
	}

	// dial
	tcon, err := tr.DialContext(ctx, ns)
	if err != nil {
		return nil, errors.New("Could not dial")
	}

	// create a new Conn that will be used to hold underlaying connections
	// from transports, middleware, as well as information about the
	// two parties.
	c := newConnWrapper(tcon, st[1:])

	return c, nil
}

func (f *Fabric) getTransport(ns string) (Transport, error) {
	// get protocol
	pr := strings.Split(ns, ":")[0]

	// check if is transport
	tr, ok := f.transports[pr]
	if !ok {
		return nil, ErrNoTransport
	}

	return tr, nil
}

func (f *Fabric) handshake(conn *conn) error {
	hs := conn.remainingStackString() // TODO Add "STREAM" + ...
	return WriteToken(conn, []byte(hs))
}

func (f *Fabric) Listen() error {
	// TODO replace with transport listens
	// TODO handle re-listening on fail
	// go func() {
	l, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		return err
	}

	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func(conn net.Conn) {
			defer func() {
				if err := conn.Close(); err != nil {
					fmt.Println("Could not close conn", err)
				}
			}()
			if err := f.handleRequest(conn); err != nil {
				fmt.Println("Listen: Could not handle request. error:", err)
			}
		}(conn)
	}
	// }()
	// return nil
}

// Handles incoming requests.
func (f *Fabric) handleRequest(tcon net.Conn) error {
	// a client initiated a connection
	fmt.Println("handleRequest: New incoming connection")

	// wrap net.Conn in Conn
	c := newConnWrapper(tcon, []string{})

	// close the connection when we're done
	defer c.Close()

	rt, err := ReadToken(tcon)
	if err != nil {
		return err
	}

	fmt.Println("Got rt:", string(rt))

	hf, ok := f.routes[string(rt)]
	if !ok {
		return ErrNoSuchMiddleware // TODO Cannot find route or something
	}

	ctx := context.Background()
	return hf(ctx, c)
}

func (f *Fabric) Next(ctx context.Context, c Conn) error {
	if c == nil {
		// TODO is this an error?
		return nil
	}

	// get next protocol
	ns := c.(*conn).popStack()
	fmt.Println("Processing", ns)

	// get protocol
	pr := strings.Split(ns, ":")[0]

	// check if is negotiator
	ng, ok := f.negotiators[pr]
	if !ok {
		return ErrNoSuchMiddleware // TODO Switch to err no negotiator
	}

	// add current part to context
	// TODO address part is a very bad name, find better one to describe address parts
	// TODO do we even still need context with values?
	mctx := context.WithValue(ctx, ContextKeyAddressPart, ns)
	// and execute them

	nc, err := ng(mctx, c)
	if err != nil {
		return err
	}

	// and move on
	return f.Next(ctx, nc)
}
