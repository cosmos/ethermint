package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ethermint/utils"

	"github.com/cosmos/cosmos-sdk/x/evm/types"
)

var _ types.QueryServer = Keeper{}

// Account implements the Query/Account gRPC method
func (q Keeper) Account(c context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if types.IsZeroAddress(req.Address); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			sdkerrors.Wrap(types.ErrZeroAddress).Error(),
		)
	}

	ctx := sdk.UnwrapSDKContext(c)
	so := keeper.GetOrNewStateObject(ctx, addr)
	balance, err := utils.MarshalBigInt(so.Balance())
	if err != nil {
		return nil, err
	}

	return &types.QueryAccountResponse{
		Balance:  balance,
		CodeHash: so.CodeHash(),
		Nonce:    so.Nonce(),
	}, nil
}
