package utils

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	emintkeys "github.com/cosmos/ethermint/keys"
	emint "github.com/cosmos/ethermint/types"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	"github.com/tendermint/tendermint/libs/cli"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// * Code from this file is a modified version of cosmos-sdk/auth/client/utils
// * to allow for using the Ethermint keybase for signing the transaction

// GenerateOrBroadcastMsgs creates a StdTx given a series of messages. If
// the provided context has generate-only enabled, the tx will only be printed
// to STDOUT in a fully offline manner. Otherwise, the tx will be signed and
// broadcasted.
func GenerateOrBroadcastETHMsgs(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return utils.PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return completeAndBroadcastETHTxCLI(txBldr, cliCtx, msgs)
}

// BroadcastETHTx Broadcasts an Ethereum Tx not wrapped in a Std Tx
func BroadcastETHTx(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, tx *evmtypes.EthereumTxMsg) error {
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	passphrase, err := emintkeys.GetPassphrase(fromName)
	if err != nil {
		return err
	}

	ethTx, err := signEthTx(txBldr.Keybase(), fromName, passphrase, tx, txBldr.ChainID())
	if err != nil {
		return err
	}

	txEncoder := txBldr.TxEncoder()

	txBytes, err := txEncoder(ethTx)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

// completeAndBroadcastETHTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func completeAndBroadcastETHTxCLI(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	passphrase, err := emintkeys.GetPassphrase(fromName)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := buildAndSign(txBldr, fromName, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a name, passphrase, and a set of messages.
func buildAndSign(bldr authtypes.TxBuilder, name, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	return sign(bldr, name, passphrase, msg)
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func sign(bldr authtypes.TxBuilder, name, passphrase string, msg authtypes.StdSignMsg) ([]byte, error) {
	sig, err := makeSignature(bldr.Keybase(), name, passphrase, msg)
	if err != nil {
		return nil, err
	}

	txEncoder := bldr.TxEncoder()

	return txEncoder(authtypes.NewStdTx(msg.Msgs, msg.Fee, []authtypes.StdSignature{sig}, msg.Memo))
}

// MakeSignature builds a StdSignature given keybase, key name, passphrase, and a StdSignMsg.
func makeSignature(keybase crkeys.Keybase, name, passphrase string,
	msg authtypes.StdSignMsg) (sig authtypes.StdSignature, err error) {
	if keybase == nil {
		keybase, err = emintkeys.NewKeyBaseFromHomeFlag()
		if err != nil {
			return
		}
	}

	var ethTx *evmtypes.EthereumTxMsg
	ethTx, ok := msg.Msgs[0].(*evmtypes.EthereumTxMsg)
	if !ok {
		return sig, fmt.Errorf("Transaction message not an Ethereum Tx")
	}

	// TODO: Move this logic to after tx is rlp decoded in keybase Sign function
	// parse the chainID from a string to a base-10 integer
	chainID, ok := new(big.Int).SetString(msg.ChainID, 10)
	if !ok {
		return sig, emint.ErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", msg.ChainID))
	}

	privKey, err := keybase.ExportPrivateKeyObject(name, passphrase)
	if err != nil {
		return
	}

	emintKey, ok := privKey.(emintcrypto.PrivKeySecp256k1)
	if !ok {
		panic(fmt.Sprintf("invalid private key type: %T", privKey))
	}

	ethTx.Sign(chainID, emintKey.ToECDSA())

	sigBytes, pubkey, err := keybase.Sign(name, passphrase, ethTx.RLPSignBytes(chainID).Bytes())
	if err != nil {
		return
	}
	return authtypes.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}

func signEthTx(keybase crkeys.Keybase, name, passphrase string,
	ethTx *evmtypes.EthereumTxMsg, chainID string) (_ *evmtypes.EthereumTxMsg, err error) {
	if keybase == nil {
		keybase, err = emintkeys.NewKeyBaseFromHomeFlag()
		if err != nil {
			return
		}
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return ethTx, emint.ErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", chainID))
	}

	privKey, err := keybase.ExportPrivateKeyObject(name, passphrase)
	if err != nil {
		return
	}

	emintKey, ok := privKey.(emintcrypto.PrivKeySecp256k1)
	if !ok {
		panic(fmt.Sprintf("invalid private key type: %T", privKey))
	}

	ethTx.Sign(intChainID, emintKey.ToECDSA())

	return ethTx, err
}

// * This context is needed because the previous GetFromFields function would initialize a
// * default keybase to lookup the address or name. The new one overrides the keybase with the
// * ethereum compatible one

// NewCLIContextWithFrom returns a new initialized CLIContext with parameters from the
// command line using Viper. It takes a key name or address and populates the FromName and
// FromAddress field accordingly.
func NewETHCLIContext() context.CLIContext {
	var nodeURI string
	var rpc rpcclient.Client

	from := viper.GetString(flags.FlagFrom)

	genOnly := viper.GetBool(flags.FlagGenerateOnly)
	fromAddress, fromName, err := getFromFields(from, genOnly)
	if err != nil {
		fmt.Printf("failed to get from fields: %v", err)
		os.Exit(1)
	}

	if !genOnly {
		nodeURI = viper.GetString(flags.FlagNode)
		if nodeURI != "" {
			rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
		}
	}

	return context.CLIContext{
		Client:        rpc,
		Output:        os.Stdout,
		NodeURI:       nodeURI,
		From:          viper.GetString(flags.FlagFrom),
		OutputFormat:  viper.GetString(cli.OutputFlag),
		Height:        viper.GetInt64(flags.FlagHeight),
		TrustNode:     viper.GetBool(flags.FlagTrustNode),
		UseLedger:     viper.GetBool(flags.FlagUseLedger),
		BroadcastMode: viper.GetString(flags.FlagBroadcastMode),
		// Verifier:      verifier,
		Simulate:     viper.GetBool(flags.FlagDryRun),
		GenerateOnly: genOnly,
		FromAddress:  fromAddress,
		FromName:     fromName,
		Indent:       viper.GetBool(flags.FlagIndentResponse),
		SkipConfirm:  viper.GetBool(flags.FlagSkipConfirmation),
	}
}

// GetFromFields returns a from account address and Keybase name given either
// an address or key name. If genOnly is true, only a valid Bech32 cosmos
// address is returned.
func getFromFields(from string, genOnly bool) (sdk.AccAddress, string, error) {
	if from == "" {
		return nil, "", nil
	}

	if genOnly {
		addr, err := sdk.AccAddressFromBech32(from)
		if err != nil {
			return nil, "", errors.Wrap(err, "must provide a valid Bech32 address for generate-only")
		}

		return addr, "", nil
	}

	keybase, err := emintkeys.NewKeyBaseFromHomeFlag()
	if err != nil {
		return nil, "", err
	}

	var info crkeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			return nil, "", err
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			return nil, "", err
		}
	}

	return info.GetAddress(), info.GetName(), nil
}
