package keeper

//// queryOrder queries the keeper with an order filtering params and returns marshalled order if found.
//func queryOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryOrderParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params: "+err.Error())
//	}
//
//	if order := keeper.GetOrder(ctx, params.Hash); order != nil {
//		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryOrderResponse{
//			Order: order,
//		})
//		if err != nil {
//			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//		}
//		return bz, nil
//	}
//
//	err := sdkerrors.Wrap(types.ErrOrderNotFound, "order does not exist: "+params.Hash.String())
//	return nil, err
//}
//
//// queryActiveOrder queries the keeper with the order filtering params and returns marshalled active order if found.
//func queryActiveOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryActiveOrderParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params:"+err.Error())
//	}
//
//	if order := keeper.GetActiveOrder(ctx, params.Hash); order != nil {
//		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryActiveOrderResponse{
//			Order: order,
//		})
//		if err != nil {
//			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//		}
//		return bz, nil
//	}
//
//	err := sdkerrors.Wrap(types.ErrOrderNotFound, "active order does not exist: "+params.Hash.String())
//	return nil, err
//}
//
//// queryArchiveOrder queries the keeper with the order filtering params and returns marshalled archive order if found.
//func queryArchiveOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryArchiveOrderParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params: "+err.Error())
//	}
//
//	if order := keeper.GetArchiveOrder(ctx, params.Hash); order != nil {
//		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryArchiveOrderResponse{
//			Order: order,
//		})
//		if err != nil {
//			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//		}
//		return bz, nil
//	}
//
//	err := sdkerrors.Wrap(types.ErrOrderNotFound, fmt.Sprintf("archive order %s does not exist", params.Hash.String()))
//	return nil, err
//}
//
