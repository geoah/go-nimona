package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
)

// blockGetCmd represents the blockGet command
var blockGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a block by its ID",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := restClient.R().Get("/blocks/" + args[0])
		if err != nil {
			return err
		}

		body := resp.Body()
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
			return nil
		}

		cmd.Println("block:")
		cmd.Println("  _id:", crypto.ID(block))
		for k, v := range block {
			cmd.Printf("  %s: %v\n", k, v)
		}
		cmd.Println("")
		return nil
	},
}

func init() {
	blockCmd.AddCommand(blockGetCmd)
}
