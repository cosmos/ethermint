# EVM Test Suite

The commands (script) will clone and run the contract tests [compound, synthetix, uniswap, ...] against the Ethermint's EVM. This test will comprehensively check the EVM against that of Ethereum's to verify that the results show the exact same behaviour. 

## Test Suites
The following are instructions on how to run the specified EVM test suite on the Ethermint network.

### Synthetix
```bash
# Start the ethermint node and expose rpc endpoint on :8545
# this can be done by running `./init.sh` in the ethermint dir

# Synthetix
git clone https://github.com/Synthetixio/synthetix.git
cd synthetix
npm install
npm install web3
npm i -g yarn
yarn add hardhat
npx hardhat compile
npx hardhat test --network development
```

### Uniswap
```bash
# Start the ethermint node and expose rpc endpoint on :8545
# this can be done by running `./init.sh` in the ethermint dir

# Uniswap v3
git clone https://github.com/Uniswap/uniswap-v3-core.git
yarn install
npx hardhat compile
```

Add ethermint network in hardhat.config.ts
```
networks: {
    ethermint: {
      url: 'http://127.0.0.1:8545',
      accounts: 'remote',
      gas: 'auto',
      gasPrice: 'auto',
      gasMultiplier: 1,
      timeout: 20000
    }
  }
```

```bash
# run the test suite
npx hardhat test --network ethermint
```

### Compound
```bash
# Start the ethermint node and expose rpc endpoint on :8545
# this can be done by running `./init.sh` in the ethermint dir

git clone https://github.com/compound-finance/compound-protocol.git
npm install
npx addle compile
```

Add ethermint network in saddle.config.ts
```
networks: {                                           
    ethermint: {
      providers: [                                      
        {env: "PROVIDER"},                              
        {http: "http://127.0.0.1:8545"}                 
      ],
      web3: {                                          
        gas: [
          {env: "GAS"},
          {default: "4600000"}
        ],
        gas_price: [
          {env: "GAS_PRICE"},
          {default: "12000000000"}
        ],
        options: {
          transactionConfirmationBlocks: 1,
          transactionBlockTimeout: 5
        }
      },
      accounts: [                                       
        {env: "ACCOUNT"},                               
        {unlocked: 0}                                
      ]
    },
  }
```

```bash
# run the test suite
npx saddle test -n ethermint
```

## Known Issues:

1. `ether.js` is not compatible with Ethermint.
Currently, Ethermint is unable to interact with the test suites and there is a slight incompatibility with ethers.js. The issue has been documented and has been in the `icebox` for some time ([https://github.com/cosmos/ethermint/issues/349](https://github.com/cosmos/ethermint/issues/349)). 

In the future, when compatible with `ethers.js`, the test suites can be run with the simple command and will be run against the Ethermint EVM. If the all of the tests pass without incident, the Ethermint EVM implementation can be deemed to produce the same output as the Ethereum EVM.

2. For `Uniswap V3`, Hardhat tests only work for hardhat test network when loading the node accounts.
3. For `Compound`, some tests are failing in their CI.