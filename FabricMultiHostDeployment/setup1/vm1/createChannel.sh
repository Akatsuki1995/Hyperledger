export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../vm4/crypto-config/ordererOrganizations/example.com/orderers/diggipet_orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_DIGGIPET_CA=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/peers/peer0.diggipet.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../../artifacts/channel/config/

export CHANNEL_NAME=firstchannel

setGlobalsForPeer0Diggipet(){
    export CORE_PEER_LOCALMSPID="DiggipetMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_DIGGIPET_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/users/Admin@diggipet.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer1Diggipet(){
    export CORE_PEER_LOCALMSPID="DiggipetMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_DIGGIPET_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/users/Admin@diggipet.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051
}

createChannel(){
    rm -rf ./channel-artifacts/*
    setGlobalsForPeer0Diggipet
    
    # Replace localhost with your orderer's vm IP address
    peer channel create -o 34.125.212.224:7050 -c $CHANNEL_NAME \
    --ordererTLSHostnameOverride diggipet_orderer.example.com \
    -f ./../../artifacts/channel/${CHANNEL_NAME}.tx --outputBlock ./channel-artifacts/${CHANNEL_NAME}.block \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
}

#createChannel

joinChannel(){
    setGlobalsForPeer0Diggipet
    peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
    
    setGlobalsForPeer1Diggipet
    peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
    
}

#joinChannel

updateAnchorPeers(){
    setGlobalsForPeer0Diggipet
    # Replace localhost with your orderer's vm IP address
    peer channel update -o 34.125.212.224:7050 --ordererTLSHostnameOverride diggipet_orderer.example.com -c $CHANNEL_NAME -f ./../../artifacts/channel/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
    
}

updateAnchorPeers

# removeOldCrypto

#createChannel
#joinChannel
#updateAnchorPeers