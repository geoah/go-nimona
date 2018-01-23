package fabric

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// RouterMiddleware is the selector middleware
type RouterMiddleware struct {
	Handlers map[string]Handler
}

// Name of the middleware
func (m *RouterMiddleware) Name() string {
	return "router"
}

// Handle is the middleware handler for the server
func (m *RouterMiddleware) Handle(ctx context.Context, c Conn) (context.Context, Conn, error) {
	// we need to negotiate what they need from us
	// read the next token, which is the request for the next middleware
	pr, err := ReadToken(c)
	if err != nil {
		return nil, nil, err
	}

	pf := strings.Split(string(pr), " ")
	if len(pf) != 2 {
		return nil, nil, errors.New("invalid router command format")
	}

	cm := pf[0]
	pm := pf[1]

	switch cm {
	case "SEL":
		return m.handleGet(ctx, c, pm)
	default:
		c.Close()
		return nil, nil, errors.New("invalid router command")
	}

}

func (m *RouterMiddleware) handleGet(ctx context.Context, c Conn, pm string) (context.Context, Conn, error) {
	addr := c.GetAddress()

	fmt.Println("Router.Handle: pm=", pm)
	fmt.Println("Router.Handle: stack=", addr.stack)

	// TODO not sure about append, might wanna cut the stack up to our index
	// and the append the new stack
	addr.stack = append(addr.stack, strings.Split(pm, "/")[1:]...)
	fmt.Println("Router.Handle: stack=", addr.stack)

	if err := WriteToken(c, []byte("ACK "+pm)); err != nil {
		return nil, nil, err
	}

	return ctx, c, nil
}

// Negotiate handles the client's side of the nimona middleware
func (m *RouterMiddleware) Negotiate(ctx context.Context, c Conn) (context.Context, Conn, error) {
	pr := c.GetAddress().RemainingString()
	fmt.Println("Router.Negotiate: pr=", pr)

	if err := WriteToken(c, []byte("SEL "+pr)); err != nil {
		return ctx, nil, err
	}

	if err := m.verifyResponse(c, "ACK "+pr); err != nil {
		return ctx, nil, err
	}

	return ctx, c, nil
}

func (m *RouterMiddleware) verifyResponse(c Conn, pr string) error {
	resp, err := ReadToken(c)
	if err != nil {
		return err
	}

	if string(resp) != pr {
		return errors.New("Invalid selector response")
	}

	return nil
}