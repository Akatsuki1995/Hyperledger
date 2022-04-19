export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../vm4/crypto-config/ordererOrganizations/example.com/orderers/diggipet_orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_VETORG_CA=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/peers/peer0.vet_org.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../../artifacts/channel/config/

export CHANNEL_NAME=firstchannel

setGlobalsForPeer0VetOrg() {
    export CORE_PEER_LOCALMSPID="VetorgMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_VETORG_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/users/Admin@vet_org.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051

}

setGlobalsForPeer1VetOrg() {
    export CORE_PEER_LOCALMSPID="VetorgMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_VETORG_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/users/Admin@vet_org.example.com/msp
    export CORE_PEER_ADDRESS=localhost:10051

}

fetchChannelBlock() {
    rm -rf ./channel-artifacts/*
    setGlobalsForPeer0VetOrg
    # Replace localhost with your orderer's vm IP address
    peer channel fetch 0 ./channel-artifacts/$CHANNEL_NAME.block -o 34.125.59.32:7050 \
        --ordererTLSHostnameOverride diggipet_orderer.example.com \
        -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
}

#fetchChannelBlock

joinChannel() {
    setGlobalsForPeer0VetOrg
    peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

    setGlobalsForPeer1VetOrg
    peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

}

#joinChannel

updateAnchorPeers() {
    setGlobalsForPeer0VetOrg
    # Replace localhost with your orderer's vm IP address
    peer channel update -o 34.125.59.32:7050 --ordererTLSHostnameOverride diggipet_orderer.example.com \
        -c $CHANNEL_NAME -f ./../../artifacts/channel/${CORE_PEER_LOCALMSPID}anchors.tx \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

}

updateAnchorPeers

# fetchChannelBlock
# joinChannel
# updateAnchorPeers
