package types

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func (gs GenesisState) Validate() error {
	// TODO: validate stuff in genesis
	return nil
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		TradePairs: []*TradePair{{
			Name:           "INJ/WETH",
			MakerAssetData: "0xf47261b000000000000000000000000021a6dbf5c87d02c6a44b5f60cfa25a8c1b4aa8a9",
			TakerAssetData: "0xf47261b000000000000000000000000047c5744863a3f34671127e8b49a86cf92b6bf789",
			Enabled:        true,
		}},

		DerivativeMarkets: []*DerivativeMarket{{
			Ticker:       "XAU/USDT",
			Oracle:       "0xe99d3570cfa174c74fa7b4118fd945665fd8e964",
			BaseCurrency: "0x1ecf86e2386d85b64a1f56aceb444633c0e778eb",
			Nonce:        "1",
			MarketId:     "0x95f2143f4f74f9f0f0f783cef79e3b6718821437455fd7f0ab6fb4d16a12319f",
			Enabled:      true,
		}},
	}
}
