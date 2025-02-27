### 1. Clone the Repositories

#### Clone the `go-ethereum` repository:

```bash
git clone https://github.com/PrikshitKumar/go-ethereum
cd go-ethereum
git checkout prikshit-zeeve
```

Run the tests:

```bash
go test -v ./internal/ethapi
go test -v ./core/txpool
```

Verify the following cases:

- **RPC - Mempool:**
  - `SendBlobTransaction`: [api_test.go#L2478](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/internal/ethapi/api_test.go#L2478)
  - `SendRawTransaction`: [api_test.go#L2543](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/internal/ethapi/api_test.go#L2543)

- **P2P Gossip - TxPool:**
  - `Add`: [txpool_test.go#L80](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/core/txpool/txpool_test.go#L80)

#### Clone the `erc20-tokens` repository:

```bash
cd ..
git clone https://github.com/PrikshitKumar/erc20-tokens
cd erc20-tokens
git checkout prikshit-zeeve
forge clean
forge build
forge test -vv ./test/ERC20.t.sol
```

### 2. Set Up Ethereum Node (`go-ethereum`)

1. Build the Geth binary:

```bash
cd ../go-ethereum
make geth
cd ./build/bin/
```

2. Clean and create data directory:

```bash
rm -rf ../../geth_data
mkdir -p ../../geth_data
```

### 3. Account Generation

- **Account 1**:

```bash
./geth account new
```

Sample output:
```
Your new key was generated
Public address: 0x5bbc51A95f92edf0346603A5C1e4F72BFa3daF11
Secret key file path: /Users/prikshitkumar/Library/Ethereum/keystore/UTC--2025-02-27T07-01-43.910391000Z--5bbc51a95f92edf0346603a5c1e4f72bfa3daf11
```

**Important:**  
- Share the **public address** with others to interact with the account.
- **Never share** the secret key; it controls access to your funds.
- **Backup** your key file to avoid losing access to funds.
- **Remember your password** to decrypt the key.

- **Account 2**: Repeat the same process for a second account.

3. Add both account addresses in `genesis.json` to make them pre-funded:  
   [genesis.json#L20](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/genesis.json#L20)

### 4. Configure Private Key Fetching

- Copy the secret key file paths and paste them into the script file with `Passwords` you set:
  - [fetchPrivateKey.sh#L11](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/fetchPrivateKey.sh#L11)
  - [fetchPrivateKey.sh#L25](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/fetchPrivateKey.sh#L25)

Make the script executable and run it after recreating the binary:

```bash
cd ../../
chmod u+x fetchPrivateKey.sh
```

### 5. Block a Specific Address

Add an address to block transactions by editing `api.go`:

- [api.go#L1445](https://github.com/PrikshitKumar/go-ethereum/blob/prikshit-zeeve/internal/ethapi/api.go#L1445)

### 6. Initialize and Run Geth

1. Initialize the blockchain with the `genesis.json`:

```bash
make geth
bash ./fetchPrivateKey.sh
cd ./build/bin/
./geth --datadir ../../geth_data init ../../genesis.json
```

2. Start the Geth node:

```bash
./geth --dev --datadir ../../geth_data --networkid 1337 --http --http.api eth,net,web3,personal --allow-insecure-unlock --http.corsdomain "*" --http.addr 0.0.0.0 --http.port 8545
```

### 7. Attach to the Geth Node

Open a new terminal and attach to the running Geth instance:

```bash
cd ./build/bin/
./geth attach ../../geth_data/geth.ipc
```

Verify the accounts have been funded:

```bash
eth.accounts
web3.fromWei(eth.getBalance("0x01dc5258c6d1443b64ee2490f5348cd7ecb05525"), "ether")
web3.fromWei(eth.getBalance("0xf107c25988dcb499c75b12cd27adbba70b72a7be"), "ether")
```

### 8. Deploy ERC20 Token Contract

Go to the `erc20-tokens` repository and deploy the contract:

```bash
forge create --rpc-url http://0.0.0.0:8545 --private-key <Deployer-Private-Key> --evm-version london src/ERC20.sol:TestToken --broadcast
```
```bash
forge create --rpc-url http://0.0.0.0:8545 --private-key 0xeac94471ef3ccdd7cf353e8f862497e6e8f4407d93879a1bf2da7854c7c40ac8 --evm-version london src/ERC20.sol:TestToken --broadcast
```

Note: If you encounter the error `contract was not deployed`, try running the command again.

Sample output:
```
Deployer: 0x01Dc5258c6D1443b64Ee2490f5348cD7ecb05525
Deployed to: 0x179da5e607EcC81935e13BAA38ee1b66b684250C <-- Contract Address
Transaction hash: 0x73a7ed540a7a1094488dc4c947cebc07170d160a7f5f5e101f2d7a2e4cca7143
```

### 9. Verify Token Minting

Check the balance of User2 before minting:

```bash
cast call <Contract-Address> "balanceOf(address)(uint256)" <User2 Address> --rpc-url 0.0.0.0:8545
```
```bash
cast call 0x179da5e607EcC81935e13BAA38ee1b66b684250C "balanceOf(address)(uint256)" 0xf107c25988dcb499c75b12cd27adbba70b72a7be --rpc-url 0.0.0.0:8545
```

Mint tokens:

```bash
cast send <Contract-Address> "mint(address,uint256)" <User2 Address> 500000000000000000000 --rpc-url 0.0.0.0:8545 --private-key <User1->Deployer Priv Key>
```
```bash
cast send 0x179da5e607EcC81935e13BAA38ee1b66b684250C "mint(address,uint256)" 0xf107c25988dcb499c75b12cd27adbba70b72a7be 500000000000000000000 --rpc-url 0.0.0.0:8545 --private-key 0xeac94471ef3ccdd7cf353e8f862497e6e8f4407d93879a1bf2da7854c7c40ac8 
```

Check User2's balance after minting:

```bash
cast call <Contract-Address> "balanceOf(address)(uint256)" <User2 Address> --rpc-url 0.0.0.0:8545
```
```bash
cast call 0x179da5e607EcC81935e13BAA38ee1b66b684250C "balanceOf(address)(uint256)" 0xf107c25988dcb499c75b12cd27adbba70b72a7be --rpc-url 0.0.0.0:8545
```

### 10. Verify Transaction Blocking

Attempt a transaction from a blocked address:

```bash
cast send <Contract-Address> "transfer(address,uint256)" <User1 Address> 100000000000000000000 --rpc-url 0.0.0.0:8545 --private-key <User2 / Blocked User Priv Key>
```
```bash
cast send 0x179da5e607EcC81935e13BAA38ee1b66b684250C "transfer(address,uint256)" 0x01dc5258c6d1443b64ee2490f5348cd7ecb05525 100000000000000000000 --rpc-url 0.0.0.0:8545 --private-key 0x2e099e159d235f201e3c0de494c43614884fbe0fd500bd1e8c63a086b48ec80b 
```

Sample output:
```
Error: server returned an error response: error code -32000: transaction from this address is blocked
```

Finally, check if the balance was deducted from User2’s account:

```bash
cast call <Contract-Address> "balanceOf(address)(uint256)" <User2 Address> --rpc-url 0.0.0.0:8545
```
```bash
cast call 0x179da5e607EcC81935e13BAA38ee1b66b684250C "balanceOf(address)(uint256)" 0xf107c25988dcb499c75b12cd27adbba70b72a7be --rpc-url 0.0.0.0:8545
```

Sample output (balance should remain unchanged):
```
500000000000000000000 [5e20]
```
