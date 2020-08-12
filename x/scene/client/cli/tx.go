package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	utils2 "github.com/cosmos/cosmos-sdk/x/scene/utils"
	"github.com/spf13/viper"
	"strconv"

	"github.com/cosmos/cosmos-sdk/x/scene/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Scene transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdGenSignAccount(cdc),
	)...)

	return txCmd
}

func GetCmdGenSignAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen-signer [amount]",
		Short:   "Generate a certain number of addresses for internal signature generation",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`$ %s tx %s gen-signer 100`, version.ClientName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			amount, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if err := utils2.GenSignAccount(amount, cliCtx, viper.GetString("lcd-home")); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("lcd-home", "", "LCD home")
	return cmd
}
