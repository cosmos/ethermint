package cli

import (
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/ethermint/x/orders/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
)

// NewTxCmd returns a root CLI command handler for certain modules/orders transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Orders admin subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterDerivativeMarketTxCmd(),
		NewSuspendDerivativeMarketTxCmd(),
		NewResumeDerivativeMarketTxCmd(),
		NewMsgRegisterSpotMarketTxCmd(),
		NewMsgSuspendSpotMarketTxCmd(),
		NewMsgResumeSpotMarketTxCmd(),
	)
	return txCmd
}

// NewRegisterDerivativeMarketTxCmd returns a CLI command handler for creating
// a MsgRegisterDerivativeMarket transaction.
func NewRegisterDerivativeMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-register [from_key_or_address] [market_name] [market_id] [oracle_address] [base_currency_address] [nonce]",
		Short: "Create and/or sign and broadcast a MsgRegisterDerivativeMarket transaction",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			marketName, err := validateMarketName(args[1])
			if err != nil {
				return err
			}

			marketId := args[2]
			oracleAddress := args[3]
			baseCurrencyAddress := args[4]
			nonce := args[5]
			msg := &types.MsgRegisterDerivativeMarket{
				Sender: cliCtx.GetFromAddress().String(),
				Market: &types.DerivativeMarket{
					Ticker:       marketName,
					Oracle:       oracleAddress,
					BaseCurrency: baseCurrencyAddress,
					Nonce:        nonce,
					MarketId:     marketId,
					Enabled:      true,
				},
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)

		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewSuspendDerivativeMarketTxCmd returns a CLI command handler for creating
// a MsgSuspendDerivativeMarket transaction.
func NewSuspendDerivativeMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-suspend [from_key_or_address] [market_id]",
		Short: "Create and/or sign and broadcast a MsgSuspendDerivativeMarket transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			marketId := args[1]
			msg := &types.MsgSuspendDerivativeMarket{
				Sender:   cliCtx.GetFromAddress().String(),
				MarketId: marketId,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)

		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewResumeDerivativeMarketCmd returns a CLI command handler for creating
// a MsgResumeDerivativeMarket transaction.
func NewResumeDerivativeMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-resume [from_key_or_address] [market_id]",
		Short: "Create and/or sign and broadcast a MsgResumeDerivativeMarket transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			marketId := args[1]
			msg := &types.MsgResumeDerivativeMarket{
				Sender:   cliCtx.GetFromAddress().String(),
				MarketId: marketId,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewMsgRegisterSpotMarketTxCmd returns a CLI command handler for creating a MsgRegisterSpotMarket transaction.
func NewMsgRegisterSpotMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-register [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgRegisterSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[3])
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterSpotMarket{
				Sender:         cliCtx.GetFromAddress().String(),
				Name:           pairName,
				MakerAssetData: args[2],
				TakerAssetData: args[3],
				Enabled:        true,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewMsgSuspendSpotMarketTxCmd returns a CLI command handler for creating a MsgSuspendSpotMarket transaction.
func NewMsgSuspendSpotMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-suspend [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgSuspendSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[3])
			if err != nil {
				return err
			}
			msg := &types.MsgSuspendSpotMarket{
				Sender:         cliCtx.GetFromAddress().String(),
				Name:           pairName,
				MakerAssetData: args[2],
				TakerAssetData: args[3],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewMsgResumeSpotMarketTxCmd returns a CLI command handler for creating a MsgResumeSpotMarket transaction.
func NewMsgResumeSpotMarketTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-suspend [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgResumeSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			_, err = validateAssetDataHex(args[3])
			if err != nil {
				return err
			}
			msg := &types.MsgResumeSpotMarket{
				Sender:         cliCtx.GetFromAddress().String(),
				Name:           pairName,
				MakerAssetData: args[2],
				TakerAssetData: args[3],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func validatePairName(name string) (string, error) {
	if parts := strings.Split(name, "/"); len(parts) != 2 ||
		len(strings.TrimSpace(parts[0])) == 0 || len(strings.TrimSpace(parts[1])) == 0 {
		err := errors.New("pair name must be in format AAA/BBB")
		return "", err
	}

	return name, nil
}

func validateMarketName(name string) (string, error) {
	if parts := strings.Split(name, "/"); len(parts) != 2 ||
		len(strings.TrimSpace(parts[0])) == 0 || len(strings.TrimSpace(parts[1])) == 0 {
		err := errors.New("market name must be in format AAA/BBB")
		return "", err
	}

	return name, nil
}

func validateAssetDataHex(data string) ([]byte, error) {
	if !strings.HasPrefix(data, "0xf47261b0") {
		err := errors.New("unsupported asset data format: missing 0xf47261b0 prefix")
		return nil, err
	} else if len(data) != 74 {
		err := errors.New("wrong addet data length, expected 74")
		return nil, err
	}

	return common.FromHex(data), nil
}
