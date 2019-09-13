package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
	emint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.EthereumTxMsg:
			return handleETHTxMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized ethermint Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle an Ethereum specific tx
func handleETHTxMsg(ctx sdk.Context, keeper Keeper, msg types.EthereumTxMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return sdk.ErrUnknownRequest("Basic validation failed").Result()
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return emint.ErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", ctx.ChainID())).Result()
	}

	// TODO: move this logic into a .From() function added to EthereumTxMsg
	chainIDMul := new(big.Int).Mul(intChainID, big.NewInt(2))
	V := new(big.Int).Sub(msg.Data.V, chainIDMul)
	V.Sub(V, big.NewInt(8))

	sigHash := msg.RLPSignBytes(intChainID)
	sender, err := types.RecoverEthSig(msg.Data.R, msg.Data.S, V, sigHash)
	if err != nil {
		// TODO: Change this error
		return sdk.ErrUnknownAddress("Unknown Sender").Result()
	}

	// Create context for evm
	context := vm.Context{
		CanTransfer: core.CanTransfer, // Looks good, but double check
		Transfer:    core.Transfer,    // Looks good, but double check
		Origin:      sender,
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        new(big.Int).SetUint64(5), // TODO: doesn't seem necessary
		Difficulty:  big.NewInt(0x30000),       // TODO: doesn't seem used in call or create
		GasLimit:    ctx.GasMeter().Limit(),
		GasPrice:    ctx.MinGasPrices().AmountOf(emint.DenomDefault).Int,
	}

	vmenv := vm.NewEVM(context, keeper.csdb, params.MainnetChainConfig, vm.Config{})

	contractCreation := msg.To() == nil
	senderRef := vm.AccountRef(sender)
	var (
		gasUsed uint64
		vmerr   error
		ret     []byte
	)

	if contractCreation {
		// TODO: Check if ctx.GasMeter().Limit() matches
		ret, _, gasUsed, vmerr = vmenv.Create(senderRef, msg.Data.Payload, ctx.GasMeter().Limit(), msg.Data.Amount)
	} else {
		// Increment the nonce for the next transaction
		keeper.csdb.SetNonce(sender, keeper.csdb.GetNonce(sender)+1)
		ret, gasUsed, vmerr = vmenv.Call(senderRef, *msg.To(), msg.Data.Payload, ctx.GasMeter().Limit(), msg.Data.Amount)
	}

	// handle errors
	if vmerr != nil {
		// TODO: Change this error to custom vm error
		return sdk.ErrUnknownRequest("VM execution error").Result()
	}

	// Refund remaining gas from tx (Check these values and ensure gas is being consumed correctly)
	ctx.GasMeter().ConsumeGas(gasUsed, "VM execution")

	// add balance for the processor of the tx (determine who rewards are being processed to)
	// TODO: Double check nothing needs to be done here

	// TODO: Remove this
	fmt.Println("VM Return: ", ret)

	return sdk.Result{}
}
