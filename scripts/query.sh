#!/bin/bash

CHANNEL_NAME=mychannel
CHAINCODE_NAME=erc20token
CHAINCODE_VERSION=1.0
CHAINCODE_PATH=github.com/erc20token

# Set the environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/UserA@org1.example.com/msp

# Query balances
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "UserA"]}'
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "User2"]}'
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "User3"]}'
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n erc20token -c '{"Args":["GetBalance", "User4"]}'