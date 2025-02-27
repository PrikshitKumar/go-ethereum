#!/bin/bash

ACCOUNT_1_PASSWORD="Account1"
ACCOUNT_2_PASSWORD="Account2"

# Extract private keys using Python
PRIVATE_KEY_1=$(python3 -c "
import json
from eth_keyfile import decode_keyfile_json

keystore_path_1 = '/Users/prikshitkumar/Library/Ethereum/keystore/UTC--2025-02-27T00-05-12.828059000Z--01dc5258c6d1443b64ee2490f5348cd7ecb05525'
password_1 = 'Account1'

with open(keystore_path_1) as f:
    keyfile_data_1 = json.load(f)

private_key_1 = decode_keyfile_json(keyfile_data_1, password_1.encode()).hex()
print(private_key_1)  # Only print the key so Bash can read it
")

PRIVATE_KEY_2=$(python3 -c "
import json
from eth_keyfile import decode_keyfile_json

keystore_path_2 = '/Users/prikshitkumar/Library/Ethereum/keystore/UTC--2025-02-27T00-05-34.716517000Z--f107c25988dcb499c75b12cd27adbba70b72a7be'
password_2 = 'Account2'

with open(keystore_path_2) as f:
    keyfile_data_2 = json.load(f)

private_key_2 = decode_keyfile_json(keyfile_data_2, password_2.encode()).hex()
print(private_key_2)  # Only print the key so Bash can read it
")

echo "Account 1 Private Key: $PRIVATE_KEY_1"
echo "Account 2 Private Key: $PRIVATE_KEY_2"

# Create temp files for private keys
TEMP_KEY_1=$(mktemp)
TEMP_KEY_2=$(mktemp)

echo $PRIVATE_KEY_1 > $TEMP_KEY_1
echo $PRIVATE_KEY_2 > $TEMP_KEY_2

# Import private keys into geth with password
echo "Importing Account 1..."
echo -e "$ACCOUNT_1_PASSWORD\n$ACCOUNT_1_PASSWORD" | ./build/bin/geth account import --datadir ./geth_data $TEMP_KEY_1

echo "Importing Account 2..."
echo -e "$ACCOUNT_2_PASSWORD\n$ACCOUNT_2_PASSWORD" | ./build/bin/geth account import --datadir ./geth_data $TEMP_KEY_2

# Remove temp files
rm -f $TEMP_KEY_1 $TEMP_KEY_2

echo "Accounts successfully imported!"