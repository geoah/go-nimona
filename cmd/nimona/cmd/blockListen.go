package cmd

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
)

// blockListenCmd represents the blockListen command
var blockListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for new incoming blocks matching a pattern",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := apiAddress + "/streams/" + args[0]
		url = strings.Replace(url, "http", "ws", 1)
		dialer := websocket.DefaultDialer
		headers := http.Header{}
		headers.Set("Content-Type", "application/cbor")
		c, _, err := dialer.Dial(url, headers)
		if err != nil {
			return err
		}

		defer c.Close()

		for {
			wsMsgType, body, err := c.ReadMessage()
			if err != nil {
				return err
			}

			if wsMsgType != 2 {
				continue
			}

			block := map[string]interface{}{}
			if err := encoding.UnmarshalInto(body, &block); err != nil {
				return err
			}

			if returnRaw {
				bs, err := json.MarshalIndent(block, "", "  ")
				if err != nil {
					return err
				}

				cmd.Println(string(bs))
				continue
			}

			cmd.Println("block:")
			cmd.Println("  _id:", crypto.ID(block))
			for k, v := range block {
				cmd.Printf("  %s: %v\n", k, v)
			}
			cmd.Println("")
		}
	},
}

func init() {
	blockCmd.AddCommand(blockListenCmd)
}
