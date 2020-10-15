package orders

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type GenesisState struct {
	Orders            []*Order            `json:"orders"`
	TradePairs        []*TradePair        `json:"tradePairs"`
	DerivativeMarkets []*DerivativeMarket `json:"derivativeMarkets"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	// TODO: validate stuff in genesis
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		TradePairs: []*TradePair{{
			Name:           "INJ/WETH",
			MakerAssetData: HexBytes(common.FromHex("0xf47261b000000000000000000000000021a6dbf5c87d02c6a44b5f60cfa25a8c1b4aa8a9")),
			TakerAssetData: HexBytes(common.FromHex("0xf47261b000000000000000000000000047c5744863a3f34671127e8b49a86cf92b6bf789")),
			Enabled:        true,
		}},

		DerivativeMarkets: []*DerivativeMarket{{
			Ticker:       "XAU/USDT",
			Oracle:       HexBytes(common.HexToAddress("0xe99d3570cfa174c74fa7b4118fd945665fd8e964").Bytes()),
			BaseCurrency: HexBytes(common.HexToAddress("0x1ecf86e2386d85b64a1f56aceb444633c0e778eb").Bytes()),
			Nonce:        BigNum(rune(1)),
			MarketID: Hash{
				Hash: common.HexToHash("0x95f2143f4f74f9f0f0f783cef79e3b6718821437455fd7f0ab6fb4d16a12319f"),
			},
			Enabled: true,
		}},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, order := range data.Orders {
		keeper.SetOrder(ctx, order)
	}
	for _, tradePair := range data.TradePairs {
		keeper.SetTradePair(ctx, tradePair)
	}
	for _, market := range data.DerivativeMarkets {
		keeper.SetDerivativeMarket(ctx, market)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Orders:            k.GetAllOrders(ctx),
		TradePairs:        k.GetAllTradePairs(ctx),
		DerivativeMarkets: k.GetAllDerivativeMarkets(ctx),
	}
}
