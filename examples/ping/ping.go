package main

import (
	"context"

	zap "go.uber.org/zap"

	fabric "github.com/nimona/go-nimona-fabric"
	logging "github.com/nimona/go-nimona-fabric/logging"
)

// Ping is our example client, it simply sends a PING string and expects a PONG
type Ping struct{}

// Name of our protocol
func (p *Ping) Name() string {
	return "ping"
}

// Negotiate will be called after all the other protocol have been processed
func (p *Ping) Negotiate(fn fabric.NegotiatorFunc) fabric.NegotiatorFunc {
	// one time scope setup area for middleware
	return func(ctx context.Context, c fabric.Conn) error {
		lgr := logging.Logger(ctx).With(
			zap.Namespace("ping"),
		)

		// close conection when done
		defer c.Close()

		if rp, ok := ctx.Value(fabric.RemoteIdentityKey{}).(string); ok {
			lgr.Info("Context contains remote id", zap.String("remote.id", rp))
		}

		// send ping
		if err := c.WriteToken([]byte("PING")); err != nil {
			lgr.Error("Could not write token", zap.Error(err))
			return err
		}

		lgr.Info("Wrote token")

		// get pong
		token, err := c.ReadToken()
		if err != nil {
			lgr.Error("Could not read token", zap.Error(err))
			return err
		}

		lgr.Info("Read token", zap.String("token", string(token)))

		return nil
	}
}

// Handle ping requests
func (p *Ping) Handle(fn fabric.HandlerFunc) fabric.HandlerFunc {
	// one time scope setup area for middleware
	return func(ctx context.Context, c fabric.Conn) error {
		lgr := logging.Logger(ctx).With(
			zap.Namespace("ping"),
		)

		lgr.Info("Handling new request")

		// close connection when done
		defer c.Close()

		if rp, ok := ctx.Value(fabric.RemoteIdentityKey{}).(string); ok {
			lgr.Info("Context contains remote id", zap.String("remote.id", rp))
		}

		// remote peer pings
		token, err := c.ReadToken()
		if err != nil {
			lgr.Error("Could not read token", zap.Error(err))
			return err
		}

		lgr.Info("Read token", zap.String("token", string(token)))

		// we pong back
		if err := c.WriteToken([]byte("PONG")); err != nil {
			lgr.Error("Could not write token", zap.Error(err))
			return err
		}

		lgr.Info("Wrote token")

		// TODO return connection as it was?
		return nil
	}
}
