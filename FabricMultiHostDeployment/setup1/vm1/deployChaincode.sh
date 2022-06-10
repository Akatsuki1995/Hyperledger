export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../vm4/crypto-config/ordererOrganizations/example.com/orderers/diggipet_orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_DIGGIPET_CA=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/peers/peer0.diggipet.example.com/tls/ca.crt
export PEER0_VETORG_CA=${PWD}/../vm2/crypto-config/peerOrganizations/vet_org.example.com/peers/peer0.vet_org.example.com/tls/ca.crt
export PEER0_BREEDERSORG_CA=${PWD}/../vm3/crypto-config/peerOrganizations/breeders_org.example.com/peers/peer0.breeders_org.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../../artifacts/channel/config/


export CHANNEL_NAME=firstchannel

setGlobalsForPeer0Diggipet() {
    export CORE_PEER_LOCALMSPID="DiggipetMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_DIGGIPET_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/users/Admin@diggipet.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer1Diggipet() {
    export CORE_PEER_LOCALMSPID="DiggipetMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_DIGGIPET_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/diggipet.example.com/users/Admin@diggipet.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051

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
CC_NAME="cc_diggipet"

packageChaincode() {
    rm -rf ${CC_NAME}.tar.gz
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode package ${CC_NAME}.tar.gz \
        --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} \
        --label ${CC_NAME}_${VERSION}
    echo "===================== Chaincode is packaged on peer0.diggipet ===================== "
}
#packageChaincode

installChaincode() {
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    echo "===================== Chaincode is installed on peer0.diggipet ===================== "

}

#installChaincode

queryInstalled() {
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo PackageID is ${PACKAGE_ID}
    echo "===================== Query installed successful on peer0.diggipet on channel ===================== "
}

#queryInstalled

approveForDiggipet() {
    setGlobalsForPeer0Diggipet
    # set -x
    # Replace localhost with your orderer's vm IP address
    peer lifecycle chaincode approveformyorg -o 34.125.59.32:7050 \
        --ordererTLSHostnameOverride diggipet_orderer.example.com --tls \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
        --collections-config ./../../artifacts/src/github.com/diggipetcc/go/collections_config.json \
        --signature-policy "OR('DiggipetMSP.member','VetOrgMSP.member','BreedersorgMSP.member')" \
        --init-required --package-id ${PACKAGE_ID} \
        --sequence 2
    # set +x

    echo "===================== chaincode approved from DiggipetOrg ===================== "

}

#queryInstalled
#approveForDiggipet

checkCommitReadyness() {
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode checkcommitreadiness \
        --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
        --collections-config /home/ubuntu/Hyperledger/FabricMultiHostDeployment/artifacts/src/github.com/diggipetcc/go/collections_config.json \
        --signature-policy "OR('DiggipetMSP.member','VetOrgMSP.member','BreedersorgMSP.member')" \
        --sequence 2 --output json --init-required
    echo "===================== checking commit readyness from DiggipetOrg ===================== "
}

#checkCommitReadyness

commitChaincodeDefinition() {
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode commit -o 34.125.59.32:7050 --ordererTLSHostnameOverride diggipet_orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_DIGGIPET_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_VETORG_CA \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_BREEDERSORG_CA \
        --version ${VERSION} --sequence ${VERSION} --init-required
}

# commitChaincodeDefinition

queryCommitted() {
    setGlobalsForPeer0Diggipet
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME}

}

#queryCommitted

chaincodeInvokeInit() {
    setGlobalsForPeer0Diggipet
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME} \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_VETORG_CA \
         --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_BREEDERSORG_CA \
        --isInit -c '{"Args":[]}'

}

 # --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_DIGGIPET_CA \

# chaincodeInvokeInit

chaincodeInvoke() {
    setGlobalsForPeer0Diggipet

    ## Create Pet
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride diggipet_orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME} \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_DIGGIPET_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_VETORG_CA   \
        --peerAddresses localhost:11051 --tlsRootCertFiles $PEER0_BREEDERSORG_CA \
        -c '{"function": "createPet","Args":["Car-Test", "lousou", "Cat", "Mixed", "Nacir"]}'

    ## Init ledger
    # peer chaincode invoke -o localhost:7050 \
    #     --ordererTLSHostnameOverride diggipet_orderer.example.com \
    #     --tls $CORE_PEER_TLS_ENABLED \
    #     --cafile $ORDERER_CA \
    #     -C $CHANNEL_NAME -n ${CC_NAME} \
    #     --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_VETORG_CA \
    #     --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_BREEDERSORG_CA \
    #     -c '{"function": "initLedger","Args":[]}'

}

# chaincodeInvoke

chaincodeQuery() {
    setGlobalsForPeer0Diggipet

    # Query Car by Id
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "queryCar","Args":["CAR0"]}'
 
}

# chaincodeQuery

# Run this function if you add any new dependency in chaincode
# presetup

# packageChaincode
# installChaincode
# queryInstalled
# approveForDiggipet
# checkCommitReadyness
# approveForMyOrg2
# checkCommitReadyness
# commitChaincodeDefination
# queryCommitted
# chaincodeInvokeInit
# sleep 5
# chaincodeInvoke
# sleep 3
# chaincodeQuery
