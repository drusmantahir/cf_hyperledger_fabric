# Chaincode Operations Automation

This repository provides a Python utility for interacting with a Hyperledger Fabric network. The `ChaincodeOperations` class simplifies the process of invoking and querying chaincode on the Fabric network, making it easier to manage blockchain operations programmatically.

## Features

- **Create Key**: Adds a key to the blockchain.
- **Create Key with Image**: Adds a key with an associated image URL.
- **Query All Keys**: Retrieves all keys stored on the blockchain.
- **Read Key**: Reads the details of a specific key.
- **Consume Key**: Marks a key as consumed or removes it.

## Prerequisites

Before using this script, ensure you have the following installed and configured:

1. **Hyperledger Fabric**:
   - Fabric network must be running.
   - Set up a channel and chaincode as per your requirements.

2. **Python**:
   - Python 3.7 or later.

3. **Fabric CLI**:
   - Ensure `peer` CLI is properly configured and available in your system's `PATH`.

4. **Certificates**:
   - TLS certificates for orderer and peers must be accessible via absolute paths.

5. **Environment Variables**:
   Export the necessary environment variables for Fabric CLI:
   ```bash
   export CORE_PEER_TLS_ENABLED=true
   export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
   export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
   
   
6. **Clone this repository:**:
   git clone https://github.com/yourusername/chaincode-operations.git
cd chaincode-operations