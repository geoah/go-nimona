package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"nimona.io/go/api"
	"nimona.io/go/dht"
	"nimona.io/go/net"
	"nimona.io/go/peers"
	"nimona.io/go/storage"
	"nimona.io/go/telemetry"
)

var (
	daemonConfigPath     string
	daemonPort           int
	daemonAPIPort        int
	daemonEnableRelaying bool
	daemonEnableMetrics  bool
	apiToken             string

	bootstrapAddresses = []string{
		"tcp:andromeda.nimona.io:21013",
		// "tcp:borealis.nimona.io:21013",
		// "tcp:cassiopeia.nimona.io:21013",
		// "tcp:draco.nimona.io:21013",
		// "tcp:eridanus.nimona.io:21013",
		// "tcp:fornax.nimona.io:21013",
		// "tcp:gemini.nimona.io:21013",
		// "tcp:hydra.nimona.io:21013",
		// "tcp:indus.nimona.io:21013",
		// "tcp:lacerta.nimona.io:21013",
		// "tcp:mensa.nimona.io:21013",
		// "tcp:norma.nimona.io:21013",
		// "tcp:orion.nimona.io:21013",
		// "tcp:pyxis.nimona.io:21013",
		// "tcp:stats.nimona.io:21013",
	}
)

// daemonStartCmd represents the daemon command
var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a peer as a daemon",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if daemonConfigPath == "" {
			usr, _ := user.Current()
			daemonConfigPath = path.Join(usr.HomeDir, ".nimona")
		}

		if err := os.MkdirAll(daemonConfigPath, 0777); err != nil {
			return errors.Wrap(err, "could not create config dir")
		}

		addressBook, err := peers.NewAddressBook(daemonConfigPath)
		if err != nil {
			return errors.Wrap(err, "could not load key")
		}

		if len(bootstrapAddresses) > 0 {
			cmd.Println("Adding bootstrap nodes")
			for _, v := range bootstrapAddresses {
				cmd.Println("  *", v)
			}
		} else {
			cmd.Println("No bootstrap nodes provided")
		}

		addressBook.LocalHostname = announceHostname

		if daemonEnableRelaying {
			if len(bootstrapAddresses) > 0 {
				cmd.Println("Relaying enabled, using bootstrap nodes")
				if err := addressBook.AddLocalPeerRelay(bootstrapAddresses...); err != nil {
					// TODO handle error
				}
			} else {
				cmd.Println("Relaying not enabled, no bootstrap nodes provided")
			}
		} else {
			cmd.Println("Relaying not enabled")
		}

		storagePath := path.Join(daemonConfigPath, "storage")
		dpr := storage.NewDiskStorage(storagePath)
		n, _ := net.NewExchange(addressBook, dpr, fmt.Sprintf("0.0.0.0:%d", daemonPort))
		dht, _ := dht.NewDHT(n, addressBook, bootstrapAddresses)
		telemetry.NewTelemetry(n, addressBook.GetLocalPeerKey(), "tcp:stats.nimona.io:21013")

		n.RegisterDiscoverer(dht)

		peerAddress := fmt.Sprintf("0.0.0.0:%d", daemonAPIPort)
		apiAddress := fmt.Sprintf("http://localhost:%d", daemonAPIPort)

		cmd.Println("Started daemon")
		cmd.Println("* Peer keys:\n  *", addressBook.GetLocalPeerInfo().Thumbprint())
		peerAddresses := addressBook.GetLocalPeerAddresses()
		cmd.Println("* Peer addresses:")
		if len(peerAddresses) > 0 {
			for _, addr := range addressBook.GetLocalPeerAddresses() {
				cmd.Println("  *", addr)
			}
		} else {
			cmd.Println("  * No addresses available")
		}
		cmd.Println("* HTTP API address:\n  *", apiAddress)

		a := api.New(addressBook, dht, n, dpr, apiToken)

		go func() {
			if err := a.Serve(peerAddress); err != nil {
				log.Fatal("Server stoping, error:", err)
			}
		}()

		api.Wait()

		return errors.Wrap(err, "http server stopped")
	},
}

func init() {
	daemon.AddCommand(daemonStartCmd)

	daemonStartCmd.PersistentFlags().IntVar(
		&daemonPort,
		"port",
		0,
		"peer port",
	)

	daemonStartCmd.PersistentFlags().IntVar(
		&daemonAPIPort,
		"api-port",
		8030,
		"api port",
	)

	daemonStartCmd.PersistentFlags().BoolVar(
		&daemonEnableRelaying,
		"relay",
		true,
		"enable relaying through bootstrap peers",
	)

	daemonStartCmd.PersistentFlags().BoolVar(
		&daemonEnableMetrics,
		"metrics",
		false,
		"enable sending anonymous metrics",
	)

	daemonStartCmd.PersistentFlags().StringSliceVar(
		&bootstrapAddresses,
		"bootstraps",
		bootstrapAddresses,
		"bootstrap addresses",
	)

	daemonStartCmd.PersistentFlags().StringVar(
		&apiToken,
		"api-token",
		apiToken,
		"api token",
	)
}
