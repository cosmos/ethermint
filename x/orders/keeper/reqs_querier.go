package keeper

//func queryZeroExTransaction(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryZeroExTransactionParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
//	}
//
//	txInfo := keeper.GetZeroExTransaction(ctx, params.TxHash)
//	if txInfo == nil {
//		return nil, nil
//	}
//	txResp := &types.QueryZeroExTransactionResponse{
//		TxType: txInfo.Type,
//	}
//
//	if txInfo.Type == types.ZeroExOrderFillRequestTx {
//		txResp.FillRequests = keeper.ListOrderFillRequestsByTxHash(ctx, params.TxHash)
//	} else if txInfo.Type == types.ZeroExOrderSoftCancelRequestTx {
//		txResp.SoftCancelRequests = keeper.ListOrderSoftCancelRequestsByTxHash(ctx, params.TxHash)
//	}
//
//	bz, err := keeper.cdc.MarshalBinaryBare(txResp)
//
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//	}
//
//	return bz, nil
//}

//func querySoftCancelledOrders(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QuerySoftCancelledOrdersParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
//	}
//
//	hashes := keeper.FindAllSoftCancelledOrders(ctx, params.OrderHashes)
//	bz, err := keeper.cdc.MarshalBinaryBare(&types.QuerySoftCancelledOrdersResponse{
//		OrderHashes: hashes,
//	})
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//	}
//
//	return bz, nil
//}

//func queryOutstandingFillRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryOutstandingFillRequestsParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
//	}
//
//	fillReqs := keeper.ListOrderFillRequestsByTxHash(ctx, params.TxHash)
//	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryOutstandingFillRequestsResponse{
//		FillRequests: fillReqs,
//	})
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//	}
//
//	return bz, nil
//}

//func queryOrderFillRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryOrderFillRequestsParams
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
//	}
//
//	fillReqs := keeper.ListOrderFillRequests(ctx, params.OrderHash)
//	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryOrderFillRequestsResponse{
//		FillRequests: fillReqs,
//	})
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//	}
//
//	return bz, nil
//}
