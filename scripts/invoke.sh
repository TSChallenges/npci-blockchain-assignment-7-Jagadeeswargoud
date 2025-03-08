#!/bin/bash

## Note: prereqs
# 1. Clone `fabric-samples` repo
# 2. Move `chaincode` directory to `fabric-samples/erc20token/chaincode/`
# 3. Move this script to `fabric-samples/test-network` directory

./network.sh down

CHANNEL_NAME=mychannel
CHAINCODE_NAME=erc20token
CHAINCODE_VERSION=1.0
CHAINCODE_PATH=github.com/erc20token


./network.sh up createChannel -ca

./network.sh deployCC -ccn erc20token -ccp ../erc20token/chaincode/ -ccl go -ccv 1.0

export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

fabric-ca-client register --id.name UserA --id.secret userapw --id.type client --id.affiliation org1 --id.attrs 'user=UserA:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

fabric-ca-client enroll -u https://UserA:userapw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/UserA@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/UserA@org1.example.com/msp/config.yaml"


fabric-ca-client register --id.name User2 --id.secret user2pw --id.type client --id.affiliation org1 --id.attrs 'user=User2:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

fabric-ca-client enroll -u https://User2:user2pw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/User2@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/User2@org1.example.com/msp/config.yaml"

fabric-ca-client register --id.name User3 --id.secret user3pw --id.type client --id.affiliation org1 --id.attrs 'user=User3:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

fabric-ca-client enroll -u https://User3:user3pw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/User3@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/User3@org1.example.com/msp/config.yaml"

fabric-ca-client register --id.name User4 --id.secret user4pw --id.type client --id.affiliation org1 --id.attrs 'user=User4:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

fabric-ca-client enroll -u https://User4:user4pw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/User4@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/User4@org1.example.com/msp/config.yaml"

# Set the environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")

# InitLedger
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/UserA@org1.example.com/msp
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"function":"InitLedger","Args":[]}'

# Invoke minting
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"function":"MintTokens","Args":["10000"]}'

# Get Admin Balance
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "UserA"]}'

# Invoke transfer tokens
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"function":"TransferTokens","Args":["UserA", "User2", "100"]}'

# Get User2 Balance
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User2@org1.example.com/msp
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "User2"]}'

# Approve Token Spend
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["ApproveSpender", "User2", "User3", "20"]}'

# TransferFrom User2 from Approved Spender User3
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User3@org1.example.com/msp
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"function":"TransferFrom","Args":["User2", "User3", "User4", "10"]}'

# GetBalance for User4
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "User4"]}'

# BurnTokens
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/UserA@org1.example.com/msp
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"function":"BurnTokens","Args":["500"]}'