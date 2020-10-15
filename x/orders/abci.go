package orders

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

)

// defaultCatchingUpOffset allows to detect if node is syncing state,
// not producing real-time blocks.
const defaultCatchingUpOffset = -10 * time.Second

func (am AppModule) BeginBlocker(ctx sdk.Context) {
	//if ctx.BlockTime().Before(time.Now().Add(defaultCatchingUpOffset)) {
	//	return
	//} else if !am.cosmosClient.CanSignTransactions() {
	//	return
	//}

	//metrics.ReportFuncCall(am.svcTags)
	//doneFn := metrics.ReportFuncTiming(am.svcTags)
	//defer doneFn()

	//evmOrdersBlockNum := am.ethOrderEventDB.CurrentBlock()
	//evmFuturesBlockNum := am.ethFuturesPositionEventDB.CurrentBlock()
	//endblockLog := log.WithFields(log.Fields{
	//	"tm_block":          ctx.BlockHeight(),
	//	"orders_evm_block":  evmOrdersBlockNum,
	//	"futures_evm_block": evmFuturesBlockNum,
	//})

	//proposerAcc := proposerAt(ctx, am.accKeeper, false)
	//if proposerAcc.Empty() {
	//	endblockLog.Infoln("no proposer for current tendermint block")
	//} else {
	//	endblockLog.Infoln("current proposer for order updates:", proposerAcc.String())
	//
	//	if am.cosmosClient.FromAddress().Equals(proposerAcc) {
	//		if am.cosmosClient == nil {
	//			endblockLog.Errorln("selected to submit order updates but loopback client is not initialized")
	//			metrics.ReportFuncError(am.svcTags)
	//			return
	//		}
	//
	//		am.proposeOrderFillUpdates(ctx, proposerAcc, evmOrdersBlockNum)
	//		am.proposeOrderCancelUpdates(ctx, proposerAcc, evmOrdersBlockNum)
	//		am.proposeFuturesPositionFillUpdates(ctx, proposerAcc, evmFuturesBlockNum)
	//		am.proposeFuturesPositionCancelUpdates(ctx, proposerAcc, evmFuturesBlockNum)
	//
	//		return
	//	}
	//}
}

const (
	defaultOnlineThreshold       = time.Minute
	defaultForgetBlocksThreshold = 5
)

//func (am AppModule) proposeOrderFillUpdates(ctx sdk.Context, proposer sdk.AccAddress, blockNum uint64) {
//	metrics.ReportFuncCall(am.svcTags)
//	doneFn := metrics.ReportFuncTiming(am.svcTags)
//	defer doneFn()
//
//	proposerLog := log.WithFields(log.Fields{
//		"tm_block":  ctx.BlockHeight(),
//		"evm_block": blockNum,
//		"proposer":  proposer.String(),
//	})
//
//	syncStatus := am.keeper.GetEvmSyncStatus(ctx)
//	if syncStatus.LatestBlockSynced > am.evmSyncStatus.LatestBlockSynced {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("updating evm sync status from keeper")
//		am.evmSyncStatus = *syncStatus
//	} else {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("using locally stored evm sync status")
//		syncStatus.LatestBlockSynced = am.evmSyncStatus.LatestBlockSynced
//	}
//
//	if syncStatus.LatestBlockSynced >= defaultForgetBlocksThreshold {
//		am.ethOrderEventDB.ForgetFillEvents(syncStatus.LatestBlockSynced - defaultForgetBlocksThreshold)
//	}
//
//	var lastBlockSynced uint64
//	err := am.ethOrderEventDB.RangeFillEvents(blockNum, func(ev *eventdb.OrderEvent) error {
//		if ev.BlockNum <= syncStatus.LatestBlockSynced {
//			return nil
//		}
//
//		proposerLog.Infoln("broadcasting MsgFilledSpotOrder")
//
//		if err := am.cosmosClient.QueueBroadcastMsg(MsgFilledSpotOrder{
//			Sender:   am.cosmosClient.FromAddress(),
//			BlockNum: ev.BlockNum,
//			TxHash: ComputeHash{
//				ComputeHash: ev.TxHash,
//			},
//			OrderHash: ComputeHash{
//				ComputeHash: ev.OrderHash,
//			},
//			AmountFilled: BigNum(ev.FillAmount.String()),
//		}); err != nil {
//			proposerLog.WithError(err).Errorln("failed to broadcast fill event msg")
//		}
//
//		if ev.BlockNum > lastBlockSynced {
//			lastBlockSynced = ev.BlockNum
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		proposerLog.WithError(err).Errorln("failed to range fill events in DB")
//		metrics.ReportFuncError(am.svcTags)
//	} else if lastBlockSynced != 0 {
//		am.evmSyncStatus.LatestBlockSynced = lastBlockSynced
//		proposerLog.WithField("lastBlockSynced", lastBlockSynced).Println("keeping locally latest block synced")
//	}
//}

//func (am AppModule) proposeOrderCancelUpdates(ctx sdk.Context, proposer sdk.AccAddress, blockNum uint64) {
//	metrics.ReportFuncCall(am.svcTags)
//	doneFn := metrics.ReportFuncTiming(am.svcTags)
//	defer doneFn()
//
//	proposerLog := log.WithFields(log.Fields{
//		"tm_block":  ctx.BlockHeight(),
//		"evm_block": blockNum,
//		"proposer":  proposer.String(),
//	})
//
//	syncStatus := am.keeper.GetEvmSyncStatus(ctx)
//	if syncStatus.LatestBlockSynced > am.evmSyncStatus.LatestBlockSynced {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("updating evm sync status from keeper")
//		am.evmSyncStatus = *syncStatus
//	} else {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("using locally stored evm sync status")
//		syncStatus.LatestBlockSynced = am.evmSyncStatus.LatestBlockSynced
//	}
//
//	if syncStatus.LatestBlockSynced > defaultForgetBlocksThreshold {
//		am.ethOrderEventDB.ForgetCancelEvents(syncStatus.LatestBlockSynced - defaultForgetBlocksThreshold)
//	}
//
//	var lastBlockSynced uint64
//	err := am.ethOrderEventDB.RangeCancelEvents(blockNum, func(ev *eventdb.OrderEvent) error {
//		if ev.BlockNum <= syncStatus.LatestBlockSynced {
//			return nil
//		}
//
//		proposerLog.Infoln("broadcasting MsgCancelledSpotOrder")
//
//		if err := am.cosmosClient.QueueBroadcastMsg(MsgCancelledSpotOrder{
//			Sender:   proposer,
//			BlockNum: ev.BlockNum,
//			TxHash: ComputeHash{
//				ComputeHash: ev.TxHash,
//			},
//			OrderHash: ComputeHash{
//				ComputeHash: ev.OrderHash,
//			},
//		}); err != nil {
//			proposerLog.WithError(err).Errorln("failed to broadcast order cancel msg")
//		}
//
//		if ev.BlockNum > lastBlockSynced {
//			lastBlockSynced = ev.BlockNum
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		proposerLog.WithError(err).Errorln("failed to range cancel events in DB")
//		metrics.ReportFuncError(am.svcTags)
//	} else if lastBlockSynced != 0 {
//		am.evmSyncStatus.LatestBlockSynced = lastBlockSynced
//		proposerLog.WithField("lastBlockSynced", lastBlockSynced).Println("keeping locally latest block synced")
//	}
//}//
//func (am AppModule) proposeFuturesPositionFillUpdates(ctx sdk.Context, proposer sdk.AccAddress, blockNum uint64) {
//	metrics.ReportFuncCall(am.svcTags)
//	doneFn := metrics.ReportFuncTiming(am.svcTags)
//	defer doneFn()
//
//	proposerLog := log.WithFields(log.Fields{
//		"tm_block":          ctx.BlockHeight(),
//		"futures_evm_block": blockNum,
//		"proposer":          proposer.String(),
//	})
//
//	syncStatus := am.keeper.GetFuturesEvmSyncStatus(ctx)
//	if syncStatus.LatestBlockSynced > am.futuresEvmSyncStatus.LatestBlockSynced {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("updating evm sync status from keeper")
//		am.futuresEvmSyncStatus = *syncStatus
//	} else {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("using locally stored evm sync status")
//		syncStatus.LatestBlockSynced = am.futuresEvmSyncStatus.LatestBlockSynced
//	}
//
//	if syncStatus.LatestBlockSynced > defaultForgetBlocksThreshold {
//		am.ethFuturesPositionEventDB.ForgetFillEvents(syncStatus.LatestBlockSynced - defaultForgetBlocksThreshold)
//	}
//
//	var lastBlockSynced uint64
//	err := am.ethFuturesPositionEventDB.RangeFillEvents(blockNum, func(ev *eventdb.FuturesPositionEvent) error {
//		if ev.BlockNum <= syncStatus.LatestBlockSynced {
//			return nil
//		}
//
//		proposerLog.Infoln("broadcasting MsgFilledDerivativeOrder")
//
//		if err := am.cosmosClient.QueueBroadcastMsg(MsgFilledDerivativeOrder{
//			Sender:   proposer,
//			BlockNum: ev.BlockNum,
//			TxHash: ComputeHash{
//				ComputeHash: ev.TxHash,
//			},
//			MakerAddress: Address{
//				Address: ev.MakerAddress,
//			},
//			MarketID: ComputeHash{
//				ComputeHash: ev.MarketID,
//			},
//			OrderHash: ComputeHash{
//				ComputeHash: ev.OrderHash,
//			},
//			PositionID:     BigNum(ev.PositionID.String()),
//			QuantityFilled: BigNum(ev.QuantityFilled.String()),
//			ContractPrice:  BigNum(ev.ContractPrice.String()),
//			IsLong:         ev.IsLong,
//		}); err != nil {
//			proposerLog.WithError(err).Errorln("failed to broadcast futures position fill event msg")
//		}
//
//		if ev.BlockNum > lastBlockSynced {
//			lastBlockSynced = ev.BlockNum
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		proposerLog.WithError(err).Errorln("failed to range futures position fill events in DB")
//		metrics.ReportFuncError(am.svcTags)
//	} else if lastBlockSynced != 0 {
//		am.futuresEvmSyncStatus.LatestBlockSynced = lastBlockSynced
//		proposerLog.WithField("lastBlockSynced", lastBlockSynced).Println("keeping locally latest block synced")
//	}
//}
//
//func (am AppModule) proposeFuturesPositionCancelUpdates(ctx sdk.Context, proposer sdk.AccAddress, blockNum uint64) {
//	metrics.ReportFuncCall(am.svcTags)
//	doneFn := metrics.ReportFuncTiming(am.svcTags)
//	defer doneFn()
//
//	proposerLog := log.WithFields(log.Fields{
//		"tm_block":          ctx.BlockHeight(),
//		"futures_evm_block": blockNum,
//		"proposer":          proposer.String(),
//	})
//
//	syncStatus := am.keeper.GetFuturesEvmSyncStatus(ctx)
//	if syncStatus.LatestBlockSynced > am.futuresEvmSyncStatus.LatestBlockSynced {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("updating evm sync status from keeper")
//		am.futuresEvmSyncStatus = *syncStatus
//	} else {
//		proposerLog.WithField("lastBlockSynced", syncStatus.LatestBlockSynced).Println("using locally stored evm sync status")
//		syncStatus.LatestBlockSynced = am.futuresEvmSyncStatus.LatestBlockSynced
//	}
//
//	if syncStatus.LatestBlockSynced > defaultForgetBlocksThreshold {
//		am.ethFuturesPositionEventDB.ForgetCancelEvents(syncStatus.LatestBlockSynced - defaultForgetBlocksThreshold)
//	}
//
//	var lastBlockSynced uint64
//	err := am.ethFuturesPositionEventDB.RangeCancelEvents(blockNum, func(ev *eventdb.FuturesPositionEvent) error {
//		if ev.BlockNum <= syncStatus.LatestBlockSynced {
//			return nil
//		}
//
//		proposerLog.Infoln("broadcasting MsgCancelledDerivativeOrder")
//
//		if err := am.cosmosClient.QueueBroadcastMsg(MsgCancelledDerivativeOrder{
//			Sender:   proposer,
//			BlockNum: ev.BlockNum,
//			TxHash: ComputeHash{
//				ComputeHash: ev.TxHash,
//			},
//			MakerAddress: Address{
//				Address: ev.MakerAddress,
//			},
//			MarketID: ComputeHash{
//				ComputeHash: ev.MarketID,
//			},
//			OrderHash: ComputeHash{
//				ComputeHash: ev.OrderHash,
//			},
//			PositionID: BigNum(ev.PositionID.String()),
//		}); err != nil {
//			proposerLog.WithError(err).Errorln("failed to broadcast futures position cancel event msg")
//		}
//
//		if ev.BlockNum > lastBlockSynced {
//			lastBlockSynced = ev.BlockNum
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		proposerLog.WithError(err).Errorln("failed to range futures position cancel events in DB")
//		metrics.ReportFuncError(am.svcTags)
//	} else if lastBlockSynced != 0 {
//		am.futuresEvmSyncStatus.LatestBlockSynced = lastBlockSynced
//		proposerLog.WithField("lastBlockSynced", lastBlockSynced).Println("keeping locally latest block synced")
//	}
//}
//
//func onlineAccountsOnly(
//	ctx sdk.Context,
//	accs []*accounts.RelayerAccount,
//	version string,
//) []string {
//	onlineAccs := make([]string, 0, len(accs))
//
//	for _, acc := range accs {
//		if !acc.IsOnline {
//			continue
//		} else if len(version) > 0 && acc.LastVersion != version {
//			continue
//		}
//
//		delta := ctx.BlockTime().Sub(time.Unix(acc.LastSeen, 0))
//		if delta > defaultOnlineThreshold {
//			continue
//		}
//
//		onlineAccs = append(onlineAccs, acc.Address.String())
//	}
//
//	return onlineAccs
//}


const DefaultVersion = "none"

//func proposerAt(ctx sdk.Context, accKeeper accounts.Keeper, past bool) sdk.AccAddress {
//	relayerAccounts := accKeeper.GetAllRelayerAccounts(ctx)
//	onlineNow := onlineAccountsOnly(ctx, relayerAccounts, DefaultVersion)
//	if len(onlineNow) == 0 {
//		log.WithFields(log.Fields{
//			"candidates": len(relayerAccounts),
//			"version":    DefaultVersion,
//		}).Infoln("empty online peer list")
//	}
//
//	// NOTE: Use this to get more weight for heavy stakers:
//	// hashring.NewWithWeights()
//	//
//	// But don't forget about Gini coefficient!
//	// Adjusted staking amount = (Relayer's staking amount)^(1/1+G)
//	ring := hashring.New(onlineNow)
//
//	randSeed := hex.EncodeToString(ctx.BlockHeader().LastBlockId.Hash)
//	accAddr, _ := ring.GetNode(randSeed)
//	proposerAcc, err := sdk.AccAddressFromBech32(accAddr)
//	if err != nil {
//		log.WithError(err).Errorln("failed to select account address from the ring")
//		return sdk.AccAddress{}
//	}
//
//	return proposerAcc
//}
