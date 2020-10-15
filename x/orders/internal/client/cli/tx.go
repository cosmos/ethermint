package cli

import (
	"bufio"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/orders transaction commands.
func NewTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Orders admin subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewRegisterDerivativeMarketTxCmd(cdc))
	txCmd.AddCommand(NewSuspendDerivativeMarketTxCmd(cdc))
	txCmd.AddCommand(NewResumeDerivativeMarketTxCmd(cdc))

	txCmd.AddCommand(NewMsgRegisterSpotMarketTxCmd(cdc))
	txCmd.AddCommand(NewMsgSuspendSpotMarketTxCmd(cdc))
	txCmd.AddCommand(NewMsgResumeSpotMarketTxCmd(cdc))

	return txCmd
}

// NewRegisterDerivativeMarketTxCmd returns a CLI command handler for creating
// a MsgRegisterDerivativeMarket transaction.
func NewRegisterDerivativeMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-register [from_key_or_address] [market_name] [market_id] [oracle_address] [base_currency_address] [nonce]",
		Short: "Create and/or sign and broadcast a MsgRegisterDerivativeMarket transaction",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			marketName, err := validateMarketName(args[1])
			if err != nil {
				return err
			}
			nonce, _ := strconv.Atoi(args[5])
			msg := types.MsgRegisterDerivativeMarket{
				Sender:       cliCtx.GetFromAddress(),
				Ticker:       marketName,
				Oracle:       common.HexToAddress(args[3]),
				BaseCurrency: common.HexToAddress(args[4]),
				Nonce:        types.NewBigNum(big.NewInt(int64(nonce))),
				MarketID: types.Hash{
					Hash: common.HexToHash(args[2]),
				},
				Enabled: true,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}

// NewSuspendDerivativeMarketTxCmd returns a CLI command handler for creating
// a MsgSuspendDerivativeMarket transaction.
func NewSuspendDerivativeMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-suspend [from_key_or_address] [market_id]",
		Short: "Create and/or sign and broadcast a MsgSuspendDerivativeMarket transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			msg := types.MsgSuspendDerivativeMarket{
				Sender: cliCtx.GetFromAddress(),
				MarketID: types.Hash{
					Hash: common.HexToHash(args[1]),
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}

// NewResumeDerivativeMarketCmd returns a CLI command handler for creating
// a MsgResumeDerivativeMarket transaction.
func NewResumeDerivativeMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derivativemarket-resume [from_key_or_address] [market_id]",
		Short: "Create and/or sign and broadcast a MsgResumeDerivativeMarket transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			msg := types.MsgResumeDerivativeMarket{
				Sender: cliCtx.GetFromAddress(),
				MarketID: types.Hash{
					Hash: common.HexToHash(args[1]),
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]
	return cmd
}

// NewMsgRegisterSpotMarketTxCmd returns a CLI command handler for creating a MsgRegisterSpotMarket transaction.
func NewMsgRegisterSpotMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-register [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgRegisterSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			makerAssetData, err := validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			takerAssetData, err := validateAssetDataHex(args[3])
			if err != nil {
				return err
			}

			msg := types.MsgRegisterSpotMarket{
				Sender:         cliCtx.GetFromAddress(),
				Name:           pairName,
				MakerAssetData: types.HexBytes(makerAssetData),
				TakerAssetData: types.HexBytes(takerAssetData),
				Enabled:        true,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// NewMsgSuspendSpotMarketTxCmd returns a CLI command handler for creating a MsgSuspendSpotMarket transaction.
func NewMsgSuspendSpotMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-suspend [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgSuspendSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			makerAssetData, err := validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			takerAssetData, err := validateAssetDataHex(args[3])
			if err != nil {
				return err
			}

			msg := types.MsgSuspendSpotMarket{
				Sender:         cliCtx.GetFromAddress(),
				Name:           pairName,
				MakerAssetData: types.HexBytes(makerAssetData),
				TakerAssetData: types.HexBytes(takerAssetData),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// NewMsgResumeSpotMarketTxCmd returns a CLI command handler for creating a MsgResumeSpotMarket transaction.
func NewMsgResumeSpotMarketTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spotmarket-suspend [from_key_or_address] [spot_market_name] [makerAssetData] [takerAssetData]",
		Short: "Create and/or sign and broadcast a MsgResumeSpotMarket transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			pairName, err := validatePairName(args[1])
			if err != nil {
				return err
			}

			makerAssetData, err := validateAssetDataHex(args[2])
			if err != nil {
				return err
			}

			takerAssetData, err := validateAssetDataHex(args[3])
			if err != nil {
				return err
			}

			msg := types.MsgResumeSpotMarket{
				Sender:         cliCtx.GetFromAddress(),
				Name:           pairName,
				MakerAssetData: types.HexBytes(makerAssetData),
				TakerAssetData: types.HexBytes(takerAssetData),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

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
