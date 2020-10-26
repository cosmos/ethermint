
KEY="mykey"
CHAINID="ethermint-1"
MONIKER="localtestnet"

# remove existing daemon and client
rm -rf ~/.ethermint*

make install

# if $KEY exists it should be deleted
ethermintd keys add $KEY --keyring-backend test

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
ethermintd init $MONIKER --chain-id $CHAINID 

# Change parameter token denominations to aphoton
cat $HOME/.ethermint/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aphoton"' > $HOME/.ethermint/config/tmp_genesis.json && mv $HOME/.ethermint/config/tmp_genesis.json $HOME/.ethermint/config/genesis.json
cat $HOME/.ethermint/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aphoton"' > $HOME/.ethermint/config/tmp_genesis.json && mv $HOME/.ethermint/config/tmp_genesis.json $HOME/.ethermint/config/genesis.json
cat $HOME/.ethermint/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aphoton"' > $HOME/.ethermint/config/tmp_genesis.json && mv $HOME/.ethermint/config/tmp_genesis.json $HOME/.ethermint/config/genesis.json
cat $HOME/.ethermint/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aphoton"' > $HOME/.ethermint/config/tmp_genesis.json && mv $HOME/.ethermint/config/tmp_genesis.json $HOME/.ethermint/config/genesis.json

# Enable faucet
# cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["faucet"]["enable_faucet"]=true' >  $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $KEY 100000000000000000000aphoton --keyring-backend test

# Sign genesis transaction
ethermintd gentx $KEY --amount=1000000000000000000aphoton --keyring-backend test --chain-id $CHAINID

# Collect genesis tx
ethermintd collect-gentxs

# echo -e '\n\ntestnet faucet enabled'
# echo -e 'to transfer tokens to your account address use:'
# echo -e "ethermintd tx faucet request 100aphoton --from $KEY\n"


# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
ethermintd start --pruning=nothing --rpc.unsafe --trace
