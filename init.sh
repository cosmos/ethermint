
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

# increase block time (?)
cat $HOME/.ethermintd/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="30000"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

if [[ $1 == "pending" ]]; then
    echo "pending mode on; block times will be set to 30s."
    # sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.ethermintd/config/config.toml
    sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.ethermintd/config/config.toml
fi

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
