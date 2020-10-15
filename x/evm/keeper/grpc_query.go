package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ethcmn "github.com/ethereum/go-ethereum/common"

	ethermint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"
)

var _ types.QueryServer = Keeper{}

// Account implements the Query/Account gRPC method
func (q Keeper) Account(c context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsZeroAddress(req.Address) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)
	so := q.GetOrNewStateObject(ctx, ethcmn.HexToAddress(req.Address))
	balance, err := ethermint.MarshalBigInt(so.Balance())
	if err != nil {
		return nil, err
	}

	return &types.QueryAccountResponse{
		Balance:  balance,
		CodeHash: so.CodeHash(),
		Nonce:    so.Nonce(),
	}, nil
}

// Balance implements the Query/Balance gRPC method
func (q Keeper) Balance(c context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsZeroAddress(req.Address) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)

	balanceInt := q.GetBalance(ctx, ethcmn.HexToAddress(req.Address))
	balance, err := ethermint.MarshalBigInt(balanceInt)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to marshal big.Int to string",
		)
	}

	return &types.QueryBalanceResponse{
		Balance: balance,
	}, nil
}

// Storage implements the Query/Storage gRPC method
func (q Keeper) Storage(c context.Context, req *types.QueryStorageRequest) (*types.QueryStorageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsZeroAddress(req.Address) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	if types.IsEmptyHash(req.Key) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrEmptyHash.Error(),
		)
	}

	if req.Height < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			sdkerrors.ErrInvalidHeight.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if req.Height != 0 {
		ctx = ctx.WithBlockHeight(req.Height)
	}

	address := ethcmn.HexToAddress(req.Address)
	key := ethcmn.HexToHash(req.Key)

	state := q.GetState(ctx, address, key)

	return &types.QueryStorageResponse{
		Value: state.String(),
	}, nil
}

// Code implements the Query/Code gRPC method
func (q Keeper) Code(c context.Context, req *types.QueryCodeRequest) (*types.QueryCodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsZeroAddress(req.Address) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrZeroAddress.Error(),
		)
	}

	if req.Height < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			sdkerrors.ErrInvalidHeight.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if req.Height != 0 {
		ctx = ctx.WithBlockHeight(req.Height)
	}

	address := ethcmn.HexToAddress(req.Address)
	code := q.GetCode(ctx, address)

	return &types.QueryCodeResponse{
		Code: code,
	}, nil
}

// TxLogs implements the Query/TxLogs gRPC method
func (q Keeper) TxLogs(c context.Context, req *types.QueryTxLogsRequest) (*types.QueryTxLogsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsEmptyHash(req.Hash) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrEmptyHash.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)

	hash := ethcmn.HexToHash(req.Hash)
	logs, err := q.GetLogs(ctx, hash)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &types.QueryTxLogsResponse{
		Logs: types.NewTransactionLogsFromEth(hash, logs).Logs,
	}, nil
}

// BlockLogs implements the Query/BlockLogs gRPC method
func (q Keeper) BlockLogs(c context.Context, req *types.QueryBlockLogsRequest) (*types.QueryBlockLogsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsEmptyHash(req.Hash) {
		return nil, status.Error(
			codes.InvalidArgument,
			types.ErrEmptyHash.Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)

	txLogs := q.GetAllTxLogs(ctx)

	return &types.QueryBlockLogsResponse{
		TxLogs: txLogs,
	}, nil
}

// BlockBloom implements the Query/BlockBloom gRPC method
func (q Keeper) BlockBloom(c context.Context, req *types.QueryBlockBloomRequest) (*types.QueryBlockBloomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Height == 0 {
		return nil, status.Error(
			codes.InvalidArgument, "block height cannot be 0",
		)
	}

	ctx := sdk.UnwrapSDKContext(c)

	height := req.Height
	if height < 0 {
		height = ctx.BlockHeight()
	}

	bloom, found := q.GetBlockBloom(ctx, height)
	if !found {
		return nil, status.Error(
			codes.NotFound, types.ErrBloomNotFound.Error(),
		)
	}

	return &types.QueryBlockBloomResponse{
		Bloom: bloom.Bytes(),
	}, nil
}

// Params implements the Query/Params gRPC method
func (q Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}
