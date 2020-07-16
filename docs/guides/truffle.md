# Truffle

Set up a truffle-ethermint local development environment. {synopsis}

#### 1. clone and build ethermint
If you don't already have ethermint installed, you will need to install ethermint:

```
git clone https://github.com/ChainSafe/ethermint
cd ethermint
make install
```

This installs  `emintd` and `emintcli` in your path.

#### 2. install truffle

```
npm i -g truffle@v5.1.31
```

#### 3. create truffle project

```
mkdir ethermint-demo
cd ethermint-demo
truffle init
```

Create `contracts/Counter.sol` containing the following contract:
```
pragma solidity ^0.5.11;

contract Counter {
  uint256 counter = 0;

  function add() public {
    counter++;
  }

  function subtract() public {
    counter--;
  }

  function getCounter() public view returns (uint256) {
    return counter;
  }
}
```

Compile the contract: `truffle compile`

Create `test/counter_test.js` containing the following truffle test:
```
const Counter = artifacts.require("Counter")

contract('Counter', accounts => {
	const from = accounts[0]
	let counter

	before(async() => {
		counter = await Counter.new()
	})

	it('should add', async() => {
		await counter.add()
		let count = await counter.getCounter()
		assert(count == 1, `count was ${count}`)
	})
})
```

#### 4. Configure truffle-config.js

Open `truffle-config.js` and uncomment the `development` section in `networks`:
```
    development: {
     host: "127.0.0.1",     // Localhost (default: none)
     port: 8545,            // Standard Ethereum port (default: none)
     network_id: "*",       // Any network (default: none)
    },
```

#### 5. Start ethermint node and deploy contract
In the `ethermint` directory, run `init.sh` to initialize and start the local ethermint node.

In another terminal, run `emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key mykey --chain-id 8 --trace` to start the RPC server.

Back in the truffle terminal, migrate the contract using `truffle migrate --network development`

You should see logs in the ethermint daemon about the transactions (one to deploy `Migrations.sol`, one to deploy `Counter.sol`)
```
I[2020-07-15|17:35:59.934] Added good transaction                       module=mempool tx=22245B935689918D332F58E82690F02073F0453D54D5944B6D64AAF1F21974E2 res="&{CheckTx:log:\"[]\" gas_wanted:6721975 }" height=3 total=1
I[2020-07-15|17:36:02.065] Executed block                               module=state height=4 validTxs=1 invalidTxs=0
I[2020-07-15|17:36:02.068] Committed state                              module=state height=4 txs=1 appHash=76BA85365F10A59FE24ADCA87544191C2D72B9FB5630466C5B71E878F9C0A111
I[2020-07-15|17:36:02.981] Added good transaction                       module=mempool tx=84516B4588CBB21E6D562A6A295F1F8876076A0CFF2EF1B0EC670AD8D8BB5425 res="&{CheckTx:log:\"[]\" gas_wanted:6721975 }" height=4 total=1
```

#### 6. Run tests

Now, you can run `truffle test --network development` to run the truffle tests using the ethermint node.

```
Using network 'development'.


Compiling your contracts...
===========================
> Everything is up to date, there is nothing to compile.



  Contract: Counter
    ✓ should add (5036ms)


  1 passing (10s)
```