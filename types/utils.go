package types

import (
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmversion "github.com/tendermint/tendermint/version"
)

// HeaderFromABCI parses an ABCI header to tendermint types Header.
func HeaderFromABCI(header abci.Header) tmtypes.Header {
	return tmtypes.Header{
		Version: tmversion.Consensus{
			Block: tmversion.Protocol(header.Version.Block),
			App:   tmversion.Protocol(header.Version.App),
		},
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,
		LastBlockID: tmtypes.BlockID{
			Hash: header.LastBlockId.Hash,
			PartsHeader: tmtypes.PartSetHeader{
				Total: int(header.LastBlockId.PartsHeader.Total),
				Hash:  header.LastBlockId.PartsHeader.Hash,
			},
		},
		LastCommitHash:     header.LastCommitHash,
		DataHash:           header.DataHash,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,
		EvidenceHash:       header.EvidenceHash,
		ProposerAddress:    header.ProposerAddress,
	}
}
