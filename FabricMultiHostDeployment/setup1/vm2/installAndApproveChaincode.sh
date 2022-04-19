export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../vm4/crypto-config/ordererOrganizations/example.com/orderers/diggipet_orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_VETORG_CA=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/peers/peer0.vet_org.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../../artifacts/channel/config/

export CHANNEL_NAME=firstchannel

setGlobalsForPeer0Vetorg() {
    export CORE_PEER_LOCALMSPID="VetorgMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_VETORG_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/users/Admin@vet_org.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051

}

setGlobalsForPeer1Vetorg() {
    export CORE_PEER_LOCALMSPID="VetorgMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_VETORG_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/vet_org.example.com/users/Admin@vet_org.example.com/msp
    export CORE_PEER_ADDRESS=localhost:10051

}

presetup() {
    echo Vendoring Go dependencies ...
    pushd ./../../artifacts/src/github.com/diggipetcc/go
    GO111MODULE=on go mod vendor
    popd
    echo Finished vendoring Go dependencies
}
#presetup

CHANNEL_NAME="firstchannel"
CC_RUNTIME_LANGUAGE="golang"
VERSION="1"
CC_SRC_PATH="./../../artifacts/src/github.com/diggipetcc/go"
CC_NAME="diggipetcc"

packageChaincode() {
    rm -rf ${CC_NAME}.tar.gz
    setGlobalsForPeer0Vetorg
    peer lifecycle chaincode package ${CC_NAME}.tar.gz \
        --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} \
        --label ${CC_NAME}_${VERSION}
    echo "===================== Chaincode is packaged on peer0.vet_org ===================== "
}
#packageChaincode

installChaincode() {
    setGlobalsForPeer0Vetorg
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.vet_org ===================== "

}

#installChaincode

queryInstalled() {
    setGlobalsForPeer0Vetorg
    peer lifecycle chaincode queryinstalled >&log.txt

    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo PackageID is ${PACKAGE_ID}
    echo "===================== Query installed successful on peer0.vet_org on channel ===================== "
}

#queryInstalled

approveForVetOrg() {
    setGlobalsForPeer0Vetorg

    # Replace localhost with your orderer's vm IP address
    peer lifecycle chaincode approveformyorg -o 34.125.59.32:7050 \
        --ordererTLSHostnameOverride diggipet_orderer.example.com --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --version ${VERSION} --init-required --package-id ${PACKAGE_ID} \
        --sequence ${VERSION}

    echo "===================== chaincode approved from Vet Org ===================== "
}
#queryInstalled
#approveForVetOrg

checkCommitReadyness() {

    setGlobalsForPeer0Vetorg
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_VETORG_CA \
        --name ${CC_NAME} --version ${VERSION} --sequence ${VERSION} --output json --init-required
    echo "===================== checking commit readyness from Vet Org ===================== "
}

checkCommitReadyness
