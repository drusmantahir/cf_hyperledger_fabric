# Hyperledger Fabric Setup and Chaincode Deployment Guide

## Prerequisites Installation and Setup

To begin, ensure you have the required tools installed:

brew install git curl
brew install go
go version
brew install node

## Clone the Fabric Samples repository and download the Fabric binaries:

git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples
./scripts/bootstrap.sh
cd fabric-samples/test-network

## Start the test network and create a channel named unicalchannel:

./network.sh up
./network.sh createChannel -c unicalchannel

## Create a directory for the chaincode, initialize a Go module, and fetch dependencies:

cd ../..  # Navigate back to fabric-samples root
mkdir -p chaincode/keygen
cd chaincode/keygen
go mod init keygen
go mod tidy
go get github.com/hyperledger/fabric-contract-api-go
go mod tidy


## Navigate back to the test-network directory and deploy the chaincode:

cd ../../test-network
./network.sh deployCC -ccn keygen -ccp ../chaincode/keygen/ -ccl go -c unicalchannel

## Set up environment variables for using the peer CLI:

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/
. ./scripts/envVar.sh
setGlobals 1

# Chaincode Commands


## Create a Key:

peer chaincode invoke \
  --waitForEvent \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -C unicalchannel \
  -n keygen \
  -c '{"function":"CreateKey","Args":["key"]}'

## Create a Key with Image:
peer chaincode invoke \
  --waitForEvent \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -C unicalchannel \
  -n keygen \
  -c '{"function":"CreateKeyWithImage","Args":["112233445566", "http://example.com/image.jpg"]}'

## Query All Keys:

peer chaincode query \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -C unicalchannel \
  -n keygen \
  -c '{"function":"QueryAllKeys","Args":[]}'

## Read a Key:

peer chaincode query \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  -C unicalchannel \
  -n keygen \
  -c '{"function":"ReadKey","Args":["112233445566"]}'

## Consume a Key:

peer chaincode invoke \
  --waitForEvent \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -C unicalchannel \
  -n keygen \
  -c '{"function":"ConsumeKey","Args":["112233445566"]}'
