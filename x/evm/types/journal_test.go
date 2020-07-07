package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ethcmn "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/crypto"
)

type JournalTestSuite struct {
	suite.Suite

	address ethcmn.Address
	journal *journal
}

func (suite *JournalTestSuite) TestJournal_append() {
	suite.journal.append(createObjectChange{
		account: &suite.address,
	})
	suite.Require().Len(suite.journal.entries, 1)
	suite.Require().Equal(1, suite.journal.dirties[0].changes)

	suite.journal.append(balanceChange{
		account: &suite.address,
		prev:    sdk.ZeroInt(),
	})

	suite.Require().Len(suite.journal.entries, 2)
	suite.Require().Equal(2, suite.journal.dirties[0].changes)
}

func (suite *JournalTestSuite) TestJournal_substractDirty() {
	suite.journal.substractDirty(suite.address)
	suite.Require().Equal(0, suite.journal.getDirty(suite.address))

	suite.journal.addDirty(suite.address)
	suite.Require().Equal(1, suite.journal.getDirty(suite.address))

	suite.journal.substractDirty(suite.address)
	suite.Require().Equal(0, suite.journal.getDirty(suite.address))

	suite.journal.substractDirty(suite.address)
	suite.Require().Equal(0, suite.journal.getDirty(suite.address))
}

func (suite *JournalTestSuite) SetupTest() {
	privkey, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	suite.address = ethcmn.BytesToAddress(privkey.PubKey().Address().Bytes())
	suite.journal = newJournal()
}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}
