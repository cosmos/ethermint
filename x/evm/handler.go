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
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("\tPanic recovered: ", r)
	// 	}
	// }()
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
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
		return emint.ErrInvalidSender(err.Error()).Result()
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
		leftOverGas uint64
		vmerr       error
		ret         []byte
	)

	if contractCreation {
		// TODO: Check if ctx.GasMeter().Limit() matches
		ret, _, leftOverGas, vmerr = vmenv.Create(senderRef, msg.Data.Payload, msg.Data.GasLimit, msg.Data.Amount)
	} else {
		// Increment the nonce for the next transaction
		keeper.SetNonce(ctx, sender, keeper.GetNonce(ctx, sender)+1)
		// fmt.Println("\tPRE BALANCE: ", keeper.GetBalance(ctx, *msg.To()))
		// fmt.Println("\tSENDER BALANCE: ", keeper.GetBalance(ctx, sender))
		// fmt.Println("\tSENDER: ", sender.Hex())
		ret, leftOverGas, vmerr = vmenv.Call(senderRef, *msg.To(), msg.Data.Payload, msg.Data.GasLimit, msg.Data.Amount)
		// fmt.Println("\tGAS REMAINING: ", leftOverGas)
		// fmt.Println("\tRECIPIENT BALANCE: ", keeper.GetBalance(ctx, *msg.To()))
		// fmt.Println("\tPOST SENDER: ", keeper.GetBalance(ctx, sender))
		// fmt.Println("\tERROR?: ", vmerr)
	}

	// fmt.Println(vmerr)

	// handle errors
	if vmerr != nil {
		return emint.ErrVMExecution(vmerr.Error()).Result()
	}

	// Refund remaining gas from tx (Check these values and ensure gas is being consumed correctly)
	ctx.GasMeter().ConsumeGas(msg.Data.GasLimit-leftOverGas, "EVM execution")

	// add balance for the processor of the tx (determine who rewards are being processed to)
	// TODO: Double check nothing needs to be done here

	// TODO: Remove this when determined return isn't needed
	fmt.Println("VM Return: ", ret)
	keeper.csdb.Finalise(true)

	// TODO: Remove commit from tx handler (should be done at end of block)
	_, err = keeper.csdb.Commit(true)
	fmt.Println(err)

	return sdk.Result{}
}
